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
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Borrower{})
	services.BaseCrudInit(Borrower{}, []Borrower{})
}

// Approve prospective borrower
func Approve(ctx *iris.Context) {

	// map the payload
	payload := make(map[string]interface{}, 0)
	if err := ctx.ReadJSON(&payload); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ktp := payload["client_ktp"].(string)
	groupID, _ := strconv.ParseUint(payload["groupId"].(string), 10, 64)

	// tanggungan, errTanggungan := json.Marshal(payload["tanggungan"])
	// if errTanggungan != nil {
	// 	fmt.Println(errTanggungan)
	// }

	// payload["tanggungan"] = tanggungan

	dataRaw, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	// get CIF with with idCardNo = ktp
	cifData := cif.Cif{}
	if ktp != "" {
		services.DBCPsql.Table("cif").Where("\"idCardNo\" = ?", ktp).Scan(&cifData)

		if cifData.IdCardNo != "" {
			// found. use existing cif
			// get borrower id
			borrower := r.RCifBorrower{}
			services.DBCPsql.Table("r_cif_borrower").Where("\"cifId\" =?", cifData.ID).Scan(&borrower)

			// reserve one loan record for this new borrower
			loan := CreateLoan(payload)
			services.DBCPsql.Table("loan").Create(&loan)

			services.DBCPsql.Table("loan_raw").Create(&loanRaw.LoanRaw{Raw: dataRaw, LoanID: loan.ID})

			rLoanBorrower := r.RLoanBorrower{
				LoanId:     loan.ID,
				BorrowerId: borrower.BorrowerId,
			}
			services.DBCPsql.Table("r_loan_borrower").Create(&rLoanBorrower)

			// which loan pricing would we like to use?
			// get the newest one
			UseProductPricing(0, loan.ID)
			go CreateRelationLoanToGroup(loan.ID, groupID)
			go CreateRelationLoanToBranch(loan.ID, groupID)
			go CreateDisbursementRecord(loan.ID, payload["disbursementDate"].(string))
		} else {
			// not found. create new CIF
			cifData = CreateCIF(payload)
			services.DBCPsql.Table("cif").Create(&cifData)

			// reserve one row for this new borrower
			newBorrower := &Borrower{Village: payload["client_desa"].(string)}
			services.DBCPsql.Table("borrower").Create(newBorrower)

			rCifBorrower := r.RCifBorrower{
				CifId:      cifData.ID,
				BorrowerId: newBorrower.ID,
			}
			services.DBCPsql.Table("r_cif_borrower").Create(&rCifBorrower)

			// reserve one loan record for this new borrower
			loan := CreateLoan(payload)
			services.DBCPsql.Table("loan").Create(&loan)

			services.DBCPsql.Table("loan_raw").Create(&loanRaw.LoanRaw{Raw: dataRaw, LoanID: loan.ID})

			rLoanBorrower := r.RLoanBorrower{
				LoanId:     loan.ID,
				BorrowerId: newBorrower.ID,
			}
			services.DBCPsql.Table("r_loan_borrower").Create(&rLoanBorrower)

			// which loan pricing would we like to use?
			// get the newest one
			UseProductPricing(0, loan.ID)
			go CreateRelationLoanToGroup(loan.ID, groupID)
			go CreateRelationLoanToBranch(loan.ID, groupID)
			go CreateDisbursementRecord(loan.ID, payload["disbursementDate"].(string))
		}

	} else {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": "Parameter is NOT valid.",
		})
	}

}

// UseProductPricing - select product pricing based on current date
func UseProductPricing(investorId uint64, loanId uint64) {
	pPricing := productPricing.ProductPricing{}
	currentDate := time.Now().Local()
	services.DBCPsql.Table("product_pricing").Where("? between \"startDate\" and \"endDate\"", currentDate).Scan(&pPricing)

	rInvProdPriceLoan := r.RInvestorProductPricingLoan{
		InvestorId:       investorId,
		ProductPricingId: pPricing.ID,
		LoanId:           loanId,
	}
	services.DBCPsql.Table("r_investor_product_pricing_loan").Create(&rInvProdPriceLoan)
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

	fmt.Printf("%+v", newCif)
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
func CreateRelationLoanToGroup(loanID uint64, groupID uint64) {
	rLoanGroupSchema := &r.RLoanGroup{LoanId: loanID, GroupId: groupID}
	services.DBCPsql.Table("r_loan_group").Create(&rLoanGroupSchema)
}

// CreateRelationLoanToBranch - create relation loan to branch
func CreateRelationLoanToBranch(loanID uint64, groupID uint64) {
	rGroupBranch := r.RGroupBranch{}
	services.DBCPsql.Table("r_group_branch").Where("\"groupId\" = ?", groupID).First(&rGroupBranch)

	rLoanBranch := &r.RLoanBranch{LoanId: loanID, BranchId: rGroupBranch.BranchId}
	services.DBCPsql.Table("r_loan_branch").Create(&rLoanBranch)
}

// CreateDisbursementRecord - Create a new disbursement record
func CreateDisbursementRecord(loanID uint64, disbursementDate string) {
	disbursementSchema := &disbursement.Disbursement{DisbursementDate: disbursementDate, Stage: "PENDING"}
	services.DBCPsql.Table("disbursement").Create(&disbursementSchema)

	rLoanDisbursementSchema := &r.RLoanDisbursement{LoanId: loanID, DisbursementId: disbursementSchema.ID}
	services.DBCPsql.Table("r_loan_disbursement").Create(&rLoanDisbursementSchema)
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
	query := "SELECT DISTINCT COUNT(borrower.id)"
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
