package loan

import (
	"fmt"

	loanHistory "bitbucket.org/go-mis/modules/loan-history"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Loan{})
	services.BaseCrudInit(Loan{}, []Loan{})
}

// FetchAll - fetchAll Loan data
func FetchAll(ctx *iris.Context) {
	loans := []LoanFetch{}

	query := "SELECT DISTINCT loan.*, "
	query += "sector.\"name\" as \"sector\", "
	query += "cif.\"name\" as \"borrower\", "
	query += "\"group\".\"name\" as \"group\", "
	query += "branch.\"name\" as \"branch\",  "
	query += "disbursement.\"disbursementDate\" "
	query += "FROM loan "
	query += "LEFT JOIN r_loan_sector ON r_loan_sector.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN sector ON r_loan_sector.\"sectorId\" = sector.\"id\" "
	query += "LEFT JOIN r_loan_borrower ON r_loan_borrower.\"borrowerId\" = loan.\"id\" "
	query += "LEFT JOIN borrower ON r_loan_borrower.\"borrowerId\" = borrower.\"id\" "
	query += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.\"id\" "
	query += "LEFT JOIN cif ON r_cif_borrower.\"cifId\" = cif.\"id\" "
	query += "LEFT JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN \"group\" ON r_loan_group.\"groupId\" = \"group\".\"id\" "
	query += "LEFT JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN branch ON r_loan_branch.\"branchId\" = branch.\"id\" "
	query += "LEFT JOIN r_loan_disbursement ON r_loan_disbursement.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN disbursement ON disbursement.\"id\" = r_loan_disbursement.\"disbursementId\" "

	services.DBCPsql.Raw(query).Find(&loans)
	ctx.JSON(iris.StatusOK, iris.Map{"data": loans})
}

// UpdateStage - Update Stage Loan
func UpdateStage(ctx *iris.Context) {
	loanData := Loan{}

	loanID := ctx.Param("id")
	stage := ctx.Param("stage")
	services.DBCPsql.First(&loanData, loanID)
	if loanData == (Loan{}) {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Can't find any loan detail.",
		})
		return
	}

	loanHistoryData := loanHistory.LoanHistory{StageFrom: loanData.Stage, StageTo: "CART", Remark: "loanId=" + fmt.Sprintf("%v", loanData.ID) + "Booked change stage"}
	services.DBCPsql.Table("loan_history").Create(&loanHistoryData)

	services.DBCPsql.Table("loan").Where("\"id\" = ?", loanData.ID).UpdateColumn("stage", stage)

	rLoanHistory := r.RLoanHistory{LoanId: loanData.ID, LoanHistoryId: loanHistoryData.ID}
	go services.DBCPsql.Table("r_loan_history").Create(&rLoanHistory)

	//stage := ctx.Param("stage")

	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status":    "success",
		"stageFrom": loanData.Stage,
		"stageTo":   stage,
	})
}
