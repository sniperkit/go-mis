package borrower

import (
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/go-mis/modules/cif"
	"bitbucket.org/go-mis/modules/disbursement"
	"bitbucket.org/go-mis/modules/loan"
	loanRaw "bitbucket.org/go-mis/modules/loan-raw"
	productPricing "bitbucket.org/go-mis/modules/product-pricing"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/survey"
	"bitbucket.org/go-mis/services"
	"github.com/jinzhu/gorm"
	"github.com/kataras/go-errors"
	iris "gopkg.in/kataras/iris.v4"
	"bitbucket.org/go-mis/config"
	"encoding/json"
	"bitbucket.org/go-mis/modules/httpClient"
)

func Init() {
		services.BaseCrudInit(Borrower{}, []Borrower{})
}

func CheckBorrowerDO(idCardNo string) (bool, error) {
	borrower := struct {
		DODate *time.Time `gorm:"column:doDate" json:"doDate"`
	}{}

	q := `select borrower."doDate"
		from cif
		join r_cif_borrower rcb on rcb."cifId" = cif.id
		join borrower on borrower.id = rcb."borrowerId"
		where cif."idCardNo" = ?`

	services.DBCPsql.Raw(q, idCardNo).Scan(&borrower)

	if borrower.DODate == nil {
		return true, nil
	}
	if time.Now().Year() >= borrower.DODate.AddDate(1, 0, 0).Year() && time.Now().YearDay() > borrower.DODate.AddDate(1, 0, 0).YearDay() {
		return true, nil
	}
	return false, nil
}

// Approve prospective borrower, sourceType: OLD/NEW
func Approve(ctx *iris.Context) {
	sourceType := ctx.Param("source-type")

	// map the payload
	payload := make(map[string]interface{}, 0)
	if err := ctx.ReadJSON(&payload); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ktp, _ := payload["client_ktp"].(string)

	d, err := CheckBorrowerDO(ktp)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if d == false {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Gagal Aprrove Borrower: Borrower berada dalam masa DO",
		})
		return
	}

	fmt.Println("Request JSON Approve Borrower", payload)
	loanID := CreateBorrowerData(ctx, payload, sourceType)

	if loanID < 1 {
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"message": "Loan " + string(loanID) + " is Created",
		},
	})
}

// CreateBorrowerData - sourceType: OLD/NEW
func CreateBorrowerData(ctx *iris.Context, payload map[string]interface{}, sourceType string) uint64 {

	groupID, _ := strconv.ParseUint(payload["groupId"].(string), 10, 64)
	sectorID, _ := strconv.ParseUint(payload["data_sector"].(string), 10, 64)

	dataRaw := payload

	db := services.DBCPsql.Begin()
	borrowerId, err := GetOrCreateBorrowerId(payload, db)
	if err != nil {
		ProcessErrorAndRollback(ctx, db, "Error Create Borrower "+err.Error())
		return 0
	}

	// reserve one loan record for this new borrower
	payload["borrowerId"] = strconv.FormatUint(borrowerId, 10)
	errLoan, loan := CreateLoan(payload)
	if errLoan != nil {
		ProcessErrorAndRollback(ctx, db, "Error Create object loan "+errLoan.Error())
		return 0
	}
	if db.Table("loan").Create(&loan).Error != nil {
		ProcessErrorAndRollback(ctx, db, "Error Create Loan")
		return 0
	}

	// check whether raw data has already exist
	existingLoanRaw := []loanRaw.LoanRaw{}
	uuid, ok := payload["uuid"]
	if !ok {
		ProcessErrorAndRollback(ctx, db, "No UUID in payload")
		return 0
	}
	uuidString, castingOk := uuid.(string)
	if !castingOk {
		ProcessErrorAndRollback(ctx, db, "Error when casting UUID")
		return 0
	}
	if err := db.Table("loan_raw").Where(`"_raw"::json->>'uuid' = ?`, uuidString).Scan(&existingLoanRaw).Error; err != nil {
		fmt.Printf(err.Error())
		ProcessErrorAndRollback(ctx, db, "Error Querying Loan Raw")
		return 0
	}

	if len(existingLoanRaw) > 0 {
		ProcessErrorAndRollback(ctx, db, "Loan has already been created")
		return 0
	}

	if db.Table("loan_raw").Create(&loanRaw.LoanRaw{Raw: dataRaw, LoanID: loan.ID}).Error != nil {
		ProcessErrorAndRollback(ctx, db, "Error Create Loan Raw")
		return 0
	}

	if db.Table("r_loan_sector").Create(&r.RLoanSector{LoanId: loan.ID, SectorId: sectorID}).Error != nil {
		ProcessErrorAndRollback(ctx, db, "Error Create Loan Sector Relation")
		return 0
	}

	rLoanBorrower := r.RLoanBorrower{
		LoanId:     loan.ID,
		BorrowerId: borrowerId,
	}
	if db.Table("r_loan_borrower").Create(&rLoanBorrower).Error != nil {
		ProcessErrorAndRollback(ctx, db, "Error Create Loan Borrower Relation")
		return 0
	}

	if UseProductPricing(0, loan.ID, db) != nil {
		ProcessErrorAndRollback(ctx, db, "Error Use Product Pricing")
		return 0
	}

	if CreateRelationLoanToGroup(loan.ID, groupID, db) != nil {
		ProcessErrorAndRollback(ctx, db, "Error Create Relation to Group")
		return 0
	}

	if CreateRelationLoanToBranch(loan.ID, groupID, db) != nil {
		ProcessErrorAndRollback(ctx, db, "Error Create Relation to Branch")
		return 0
	}

	if CreateDisbursementRecord(loan.ID, payload["disbursementDate"].(string), db) != nil {
		ProcessErrorAndRollback(ctx, db, "Error Create Disbusrement")
		return 0
	}

	if sourceType == "OLD" {
		dbSurvey := services.DBCPsqlSurvey.Begin()

		idCardNo := payload["client_ktp"].(string)
		if setOldSurveyStatus(idCardNo, "APPROVE", dbSurvey) != nil {
			ProcessErrorAndRollback(ctx, db, "Error Setting Data in DB Survey")
			return 0
		}

		dbSurvey.Commit()
	} else {
		uuid := payload["uuid"].(string)
		if setNewSurveyStatus(uuid, "APPROVE", db) != nil {
			ProcessErrorAndRollback(ctx, db, "Error Setting Data in DB Survey")
			return 0
		}
	}

	db.Commit()

	return loan.ID
}

// status: APPROVE/REJECT
func setOldSurveyStatus(idCardNo string, status string, db *gorm.DB) error {
	query := "select * from a_fields where key = 'client_ktp' and val = ? order by id desc limit 1"
	aFields := survey.AFields{}
	services.DBCPsqlSurvey.Raw(query, idCardNo).Scan(&aFields)

	var statusBoolean bool
	if status == "APPROVE" {
		statusBoolean = true
	} else {
		statusBoolean = false
	}

	if err := db.Table("a_fields").Where("answer_id = ?", aFields.AnswerId).UpdateColumn("is_migrated", true).Error; err != nil {
		return err
	}

	if err := db.Table("a_fields").Where("answer_id = ?", aFields.AnswerId).UpdateColumn("is_approve", statusBoolean).Error; err != nil {
		return err
	}

	if err := db.Table("a_fields").Where("answer_id = ?", aFields.AnswerId).UpdateColumn("updated_at", "now").Error; err != nil {
		return err
	}

	return nil
}

// status: APPROVE/REJECT
func setNewSurveyStatus(uuid string, status string, db *gorm.DB) error {
	var statusBoolean bool
	if status == "APPROVE" {
		statusBoolean = true
	} else {
		statusBoolean = false
	}

	if err := db.Table("survey").Where("uuid = ? AND \"deletedAt\" IS NULL", uuid).UpdateColumn("isMigrate", true).Error; err != nil {
		return err
	}

	if err := db.Table("survey").Where("uuid = ? AND \"deletedAt\" IS NULL", uuid).UpdateColumn("isApprove", statusBoolean).Error; err != nil {
		return err
	}

	if err := db.Table("survey").Where("uuid = ? AND \"deletedAt\" IS NULL", uuid).UpdateColumn("updatedAt", "now").Error; err != nil {
		return err
	}

	return nil
}

func GetOrCreateBorrowerId(payload map[string]interface{}, db *gorm.DB) (uint64, error) {
	cifData := cif.InsertCif{}
	ktp := payload["client_ktp"].(string)
	db.Table("cif").Where("\"idCardNo\" = ?", ktp).Scan(&cifData)

	if cifData.IdCardNo != "" {
		// get
		borrower := r.RCifBorrower{}
		err := db.Table("r_cif_borrower").Where("\"cifId\" =?", cifData.ID).Scan(&borrower).Error
		if err != nil {
			return 0, err
		}
		return borrower.BorrowerId, nil
	} else {
		// create

		// create the CIF
		cifData = CreateCIF(payload)
		err := db.Table("cif").Create(&cifData).Error
		if err != nil {
			return 0, err
		}

		// create the Borrower
		newBorrower := &Borrower{Village: payload["client_desa"].(string)}
		err = db.Table("borrower").Create(newBorrower).Error
		if err != nil {
			return 0, err
		}

		// create the relation between Borrower and Cif
		rCifBorrower := r.RCifBorrower{
			CifId:      cifData.ID,
			BorrowerId: newBorrower.ID,
		}
		err = db.Table("r_cif_borrower").Create(&rCifBorrower).Error
		if err != nil {
			return 0, err
		}
		return newBorrower.ID, nil
	}
}

// UseProductPricing - select product pricing based on current date
func UseProductPricing(investorId uint64, loanId uint64, db *gorm.DB) error {
	pPricing := productPricing.ProductPricing{}
	if err := db.Table("product_pricing").Where("current_date::date between \"startDate\"::date and \"endDate\"::date and \"isInstitutional\"=false and \"deletedAt\" IS NULL").Scan(&pPricing).Error; err != nil {
		return err
	}

	rInvProdPriceLoan := r.RInvestorProductPricingLoan{
		InvestorId:       investorId,
		ProductPricingId: pPricing.ID,
		LoanId:           loanId,
	}

	if err := db.Table("r_investor_product_pricing_loan").Create(&rInvProdPriceLoan).Error; err != nil {
		print(err)
		return err
	}
	return nil
}

// CreateCIF - create CIF object
func CreateCIF(payload map[string]interface{}) cif.InsertCif {
	// convert each payload  field into empty string
	var cpl map[string]string
	cpl = make(map[string]string)
	for k, v := range payload {
		if v == nil {
			cpl[k] = ""
		} else {
			cpl[k] = v.(string)
		}
	}

	// map each payload field to it's respective cif field
	newCif := cif.InsertCif{}

	wifeIncome, _ := strconv.ParseFloat(cpl["data_pendapatan_istri"], 64)
	husbandIncome, _ := strconv.ParseFloat(cpl["data_pendapatan_suami"], 64)
	otherIncome, _ := strconv.ParseFloat(cpl["data_pendapatan_lain"], 64)

	// here the only payload that match the cif fields
	newCif.Username = cpl["client_simplename"]
	newCif.Name = cpl["client_fullname"]
	newCif.PlaceOfBirth = cpl["client_birthplace"]
	newCif.DateOfBirth = cpl["client_birthdate"]
	newCif.IdCardNo = cpl["client_ktp"]
	newCif.IdCardFilename = cpl["photo_ktp"]
	newCif.TaxCardNo = cpl["client_npwp"]
	newCif.MaritalStatus = cpl["client_marital_status"]
	newCif.MotherName = cpl["client_ibu_kandung"]
	newCif.Religion = cpl["client_religion"]
	newCif.Address = cpl["client_alamat"]
	newCif.Kelurahan = cpl["client_desa"]
	newCif.Kecamatan = cpl["client_kecamatan"]
	newCif.RT = cpl["client_rt"]
	newCif.RW = cpl["client_rw"]
	newCif.Income = wifeIncome + husbandIncome + otherIncome
	newCif.Occupation = cpl["client_job"]

	return newCif
}

// CreateLoan - create loan object
func CreateLoan(payload map[string]interface{}) (error, loan.Loan) {
	// convert each payload  field into empty string
	var cpl map[string]string
	cpl = make(map[string]string)
	for k, v := range payload {
		if v == nil {
			cpl[k] = ""
		} else {
			cpl[k] = v.(string)
		}
	}
	if cpl["tenor"] == "" ||
		cpl["plafond"] == "" ||
		cpl["rate"] == "" ||
		cpl["creditScoreGrade"] == "" ||
		cpl["creditScoreValue"] == "" {
		return errors.New("CSTrip is required"), loan.Loan{}
	}
	emptyPenanggungJawab := cpl["client_ktp_penanggung_jawab"] == "" || cpl["photo_ktp_penanggung_jawab"] == "" || cpl["photo_penanggung_jawab"] == ""

	if cpl["loanType"] != "AVARA" && emptyPenanggungJawab {
		return errors.New("Penanggung Jawab is required"), loan.Loan{}
	}
	// map each payload field to it's respective cif field
	newLoan := loan.Loan{}
	newLoan.URLPic1 = cpl["photo_client"]
	newLoan.URLPic2 = cpl["photo_client_square"]
	newLoan.Purpose = cpl["data_tujuan"]
	newLoan.SubmittedLoanDate, _ = cpl["data_tgl"] // time.Parse("2006-01-02 15:04:05", cpl["data_tgl"])
	newLoan.SubmittedTenor, _ = strconv.ParseInt(cpl["data_jangkawaktu"], 10, 64)
	newLoan.AgreementType = cpl["data_akad"]
	newLoan.SubmittedPlafond, _ = strconv.ParseFloat(cpl["data_pengajuan"], 64)
	newLoan.SubmittedInstallment, _ = strconv.ParseFloat(cpl["data_rencana_angsuran"], 64)
	newLoan.LoanPeriod, _ = strconv.ParseInt(cpl["data_ke"], 10, 64)

	newLoan.Tenor, _ = strconv.ParseUint(cpl["tenor"], 10, 64)
	newLoan.Rate, _ = strconv.ParseFloat(cpl["rate"], 64) // temporary value until the input defined in the future
	newLoan.Installment, _ = strconv.ParseFloat(cpl["installment"], 64)
	newLoan.Plafond, _ = strconv.ParseFloat(cpl["plafond"], 64)

	newLoan.CreditScoreGrade = cpl["creditScoreGrade"]
	newLoan.CreditScoreValue, _ = strconv.ParseFloat(cpl["creditScoreValue"], 64)
	newLoan.Stage = "PRIVATE"

	// set loan type if exist
	// otherwise, NORMAL will be set as default
	if cpl["loanType"] != "" {
		newLoan.LoanType = cpl["loanType"]
	}

	borrowerId := payload["borrowerId"]

	query := `SELECT loan.* FROM borrower JOIN r_loan_borrower ON r_loan_borrower."borrowerId" = borrower."id" JOIN loan ON loan."id" = r_loan_borrower."loanId" WHERE borrower."id" = ?`

	oldLoan := loan.Loan{}

	services.DBCPsql.Raw(query, borrowerId).Scan(&oldLoan)

	if (oldLoan == loan.Loan{}) {
		newLoan.IsLWK = false
		newLoan.IsUPK = false
	} else {
		newLoan.IsLWK = true
		newLoan.IsUPK = true
		newLoan.Subgroup = oldLoan.Subgroup
	}

	return nil, newLoan
}

// CreateRelationLoanToGroup - Create relation loan to group
func CreateRelationLoanToGroup(loanID uint64, groupID uint64, db *gorm.DB) error {
	rLoanGroupSchema := &r.RLoanGroup{LoanId: loanID, GroupId: groupID}
	if err := db.Table("r_loan_group").Create(&rLoanGroupSchema).Error; err != nil {
		return err
	}
	return nil
}

// CreateRelationLoanToBranch - create relation loan to branch
func CreateRelationLoanToBranch(loanID uint64, groupID uint64, db *gorm.DB) error {
	rGroupBranch := r.RGroupBranch{}

	if err := db.Table("r_group_branch").Where("\"groupId\" = ?", groupID).First(&rGroupBranch).Error; err != nil {
		return err
	}

	rLoanBranch := &r.RLoanBranch{LoanId: loanID, BranchId: rGroupBranch.BranchId}
	if err := db.Table("r_loan_branch").Create(&rLoanBranch).Error; err != nil {
		return err
	}
	return nil
}

// CreateDisbursementRecord - Create a new disbursement record
func CreateDisbursementRecord(loanID uint64, disbursementDate string, db *gorm.DB) error {
	disbursementSchema := &disbursement.Disbursement{DisbursementDate: disbursementDate, Stage: "PENDING"}
	if err := db.Table("disbursement").Create(&disbursementSchema).Error; err != nil {
		return err
	}

	rLoanDisbursementSchema := &r.RLoanDisbursement{LoanId: loanID, DisbursementId: disbursementSchema.ID}
	if err := db.Table("r_loan_disbursement").Create(&rLoanDisbursementSchema).Error; err != nil {
		return err
	}
	return nil
}

// ProspectiveBorrowerUpdateStatus - update status
func ProspectiveBorrowerUpdateStatus(ctx *iris.Context) {
	answerID := ctx.Param("id")
	services.DBCPsqlSurvey.Table("a_fields").Where("answer_id = ?", answerID).UpdateColumn("is_migrated", true)
	services.DBCPsqlSurvey.Table("a_fields").Where("answer_id = ?", answerID).UpdateColumn("is_approve", true)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}

// ProspectiveBorrowerUpdateStatusToReject - update status to `reject`
func ProspectiveBorrowerUpdateStatusToReject(ctx *iris.Context) {
	answerID := ctx.Param("id")
	services.DBCPsqlSurvey.Table("a_fields").Where("answer_id = ?", answerID).UpdateColumn("is_migrated", true)
	services.DBCPsqlSurvey.Table("a_fields").Where("answer_id = ?", answerID).UpdateColumn("is_approve", false)
	services.DBCPsqlSurvey.Table("a_fields").Where("answer_id = ?", answerID).UpdateColumn("updated_at", time.Now())

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}

type Count struct {
	Total int64 `gorm:"column:count"`
}

// GetTotalBorrowerByBranchID - get total borrower
func GetTotalBorrowerByBranchID(ctx *iris.Context) {
	query := "SELECT COUNT(DISTINCT borrower.id)"
	query += "FROM borrower "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"borrowerId\" = borrower.id "
	query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = r_loan_borrower.\"loanId\" "
	query += "WHERE r_loan_branch.\"branchId\" = ? AND borrower.\"deletedAt\" IS NULL "

	countSchema := Count{}
	services.DBCPsql.Raw(query, ctx.Param("branch_id")).Scan(&countSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   countSchema.Total,
	})
}

func ProcessErrorAndRollback(ctx *iris.Context, db *gorm.DB, message string) {
	db.Rollback()
	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status":  "error",
		"message": message,
	})
}

func GetBorrowerByGroup(ctx *iris.Context) {

	type BorrowerGroup struct {
		ID   uint64 `json:_id`
		Name string `json:name`
	}
	m := []BorrowerGroup{}

	groupId := ctx.Get("groupId")

	query := `SELECT DISTINCT (cif."name") AS "name", borrower."id" as "id" FROM "group"
	LEFT JOIN r_loan_group rlg ON rlg."groupId" = "group"."id"
	LEFT JOIN loan ON loan."id" = rlg."loanId"
	LEFT JOIN r_loan_borrower rlb ON rlb."loanId" = "loan"."id"
	LEFT JOIN borrower ON borrower."id" = rlb."borrowerId"
	LEFT JOIN r_cif_borrower rcb ON rcb."borrowerId" = "borrower"."id"
	LEFT JOIN cif ON cif."id" = rcb."cifId"
	WHERE "group"."id" = ?`

	if e := services.DBCPsql.Raw(query, groupId).Find(&m).Error; e != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "failed",
			"data":   e,
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   m,
	})

}

func GetProspectiveAvaraBorrowerByBranch(ctx *iris.Context) {
	branchID := ctx.Param("branch_id")

	goBorrowerEndpoint := fmt.Sprintf(`%s/borrower/prospective-avara/%v`, config.GoBorrowerPath, branchID)
	var responseBody struct {
		Status  int                         `json:"status"`
		Code    int                         `json:"code"`
		Message string                      `json:"message"`
		Data    []ProspectiveAvaraBorrower  `json:"data"`
	}
	resBody, err := httpClient.Get(goBorrowerEndpoint)
	if err != nil {
		fmt.Printf("error contacting go borrower: %+v", err)
		ctx.JSON(iris.StatusInternalServerError, nil)
		return
	}

	err = json.Unmarshal(resBody, &responseBody)
	if err != nil {
		fmt.Printf("error parsing response body: %+v", err)
		ctx.JSON(iris.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": responseBody.Data,
	})
}

func SubmitAvaraOffer(ctx *iris.Context) {
	payload := BorrowerAvaraRequest{}
	if err := ctx.ReadJSON(&payload); err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	goBorrowerEndpoint := fmt.Sprintf(`%s/borrower/process-avara`, config.GoBorrowerPath)
	var responseBody struct {
		Status  int                         `json:"status"`
		Code    int                         `json:"code"`
		Message string                      `json:"message"`
		Error   string                      `json:"error"`
	}
	resBody, err := httpClient.Post(goBorrowerEndpoint, payload)
	if err != nil {
		fmt.Printf("error contacting go borrower: %+v", err)
		ctx.JSON(iris.StatusInternalServerError, nil)
		return
	}

	err = json.Unmarshal(resBody, &responseBody)
	if err != nil {
		fmt.Printf("error parsing response body: %+v", err)
		ctx.JSON(iris.StatusInternalServerError, nil)
		return
	}

	if responseBody.Code != 200 {
		fmt.Printf("go-borrower response: %+v\n", responseBody)
		ctx.JSON(iris.StatusBadRequest, responseBody.Error)
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": nil,
	})
	return
}