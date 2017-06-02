package borrower

import (
	"encoding/json"
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
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Borrower{})
	services.BaseCrudInit(Borrower{}, []Borrower{})
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

	loanID := CreateBorrowerData(ctx, payload, sourceType)

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

	dataRaw, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	db := services.DBCPsql.Begin()
	borrowerId, err := GetOrCreateBorrowerId(payload, db)
	if err != nil {
		processErrorAndRollback(ctx, db, "Error Create Borrower "+err.Error())
		return 0
	}

	// reserve one loan record for this new borrower
	loan := CreateLoan(payload)
	if db.Table("loan").Create(&loan).Error != nil {
		processErrorAndRollback(ctx, db, "Error Create Loan")
		return 0
	}

	if db.Table("loan_raw").Create(&loanRaw.LoanRaw{Raw: dataRaw, LoanID: loan.ID}).Error != nil {
		processErrorAndRollback(ctx, db, "Error Create Loan Raw")
		return 0
	}

	if db.Table("r_loan_sector").Create(&r.RLoanSector{LoanId: loan.ID, SectorId: sectorID}).Error != nil {
		processErrorAndRollback(ctx, db, "Error Create Loan Sector Relation")
		return 0
	}

	rLoanBorrower := r.RLoanBorrower{
		LoanId:     loan.ID,
		BorrowerId: borrowerId,
	}
	if db.Table("r_loan_borrower").Create(&rLoanBorrower).Error != nil {
		processErrorAndRollback(ctx, db, "Error Create Loan Borrower Relation")
		return 0
	}

	if UseProductPricing(0, loan.ID, db) != nil {
		processErrorAndRollback(ctx, db, "Error Use Product Pricing")
		return 0
	}

	if CreateRelationLoanToGroup(loan.ID, groupID, db) != nil {
		processErrorAndRollback(ctx, db, "Error Create Relation to Group")
		return 0
	}

	if CreateRelationLoanToBranch(loan.ID, groupID, db) != nil {
		processErrorAndRollback(ctx, db, "Error Create Relation to Branch")
		return 0
	}

	if CreateDisbursementRecord(loan.ID, payload["disbursementDate"].(string), db) != nil {
		processErrorAndRollback(ctx, db, "Error Create Disbusrement")
		return 0
	}

	if sourceType == "OLD" {
		dbSurvey := services.DBCPsqlSurvey.Begin()

		idCardNo := payload["client_ktp"].(string)
		if setOldSurveyStatus(idCardNo, "APPROVE", dbSurvey) != nil {
			processErrorAndRollback(ctx, db, "Error Setting Data in DB Survey")
			return 0
		}

		dbSurvey.Commit()
	} else {
		uuid := payload["uuid"].(string)
		if setNewSurveyStatus(uuid, "APPROVE", db) != nil {
			processErrorAndRollback(ctx, db, "Error Setting Data in DB Survey")
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
	cifData := cif.Cif{}
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
	currentDate := time.Now().Local()
	if err := db.Table("product_pricing").Where("? between \"startDate\" and \"endDate\" and \"isInstitutional\"=false and \"deletedAt\" IS NULL", currentDate).Scan(&pPricing).Error; err != nil {
		return err
	}

	rInvProdPriceLoan := r.RInvestorProductPricingLoan{
		InvestorId:       investorId,
		ProductPricingId: pPricing.ID,
		LoanId:           loanId,
	}
	if err := db.Table("r_investor_product_pricing_loan").Create(&rInvProdPriceLoan).Error; err != nil {
		return err
	}
	return nil
}

// CreateCIF - create CIF object
func CreateCIF(payload map[string]interface{}) cif.Cif {
	// convert each payload  field into empty string
	var cpl map[string]string
	cpl = make(map[string]string)
	for k, v := range payload {
		cpl[k] = v.(string)
	}

	// map each payload field to it's respective cif field
	newCif := cif.Cif{}

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
	newCif.Income, _ = strconv.ParseFloat(cpl["data_pendapatan_istri"], 64)
	newCif.Occupation = cpl["client_job"]

	return newCif
}

// CreateLoan - create loan object
func CreateLoan(payload map[string]interface{}) loan.Loan {
	// convert each payload  field into empty string
	var cpl map[string]string
	cpl = make(map[string]string)
	for k, v := range payload {
		cpl[k] = v.(string)
	}

	// map each payload field to it's respective cif field
	newLoan := loan.Loan{}
	newLoan.URLPic1 = cpl["photo_client"]
	newLoan.URLPic2 = cpl["photo_client_square"]
	newLoan.Purpose = cpl["data_tujuan"]
	newLoan.SubmittedLoanDate, _ = cpl["data_tgl"] // time.Parse("2006-01-02 15:04:05", cpl["data_tgl"])
	newLoan.SubmittedTenor, _ = strconv.ParseInt(cpl["data_jangkawaktu"], 10, 64)
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
	newLoan.IsLWK = false
	newLoan.IsUPK = false

	return newLoan
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

func processErrorAndRollback(ctx *iris.Context, db *gorm.DB, message string) {
	db.Rollback()
	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status":  "error",
		"message": message,
	})
}
