package disbursement

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Disbursement{})
	services.BaseCrudInit(Disbursement{}, []Disbursement{})
}

// FetchAll - fetchAll Disbursement data
func FetchAll(ctx *iris.Context) {
	disbursements := []DisbursementFetch{}

	query := "SELECT disbursement.\"id\", disbursement.\"disbursementDate\", disbursement.\"stage\", "
	query += "loan.\"submittedLoanDate\", loan.\"plafond\", \"group\".\"name\" as \"group\", branch.\"name\" as \"branch\", "
	query += "cif.\"name\" as \"borrower\" "
	query += "FROM disbursement "
	query += "JOIN r_loan_disbursement ON r_loan_disbursement.\"disbursementId\" = disbursement.\"id\" "
	query += "JOIN loan ON r_loan_disbursement.\"loanId\" = loan.\"id\" "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = loan.\"id\" "
	query += "JOIN borrower ON borrower.\"id\" = r_loan_borrower.\"borrowerId\" "
	query += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.\"id\" "
	query += "JOIN cif ON cif.\"id\" = r_cif_borrower.\"cifId\" "
	query += "LEFT JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN branch ON r_loan_branch.\"branchId\" = branch.\"id\" "
	query += "LEFT JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	query += "LEFT JOIN \"group\" ON r_loan_group.\"groupId\" = \"group\".\"id\" "
	query += "WHERE disbursement.\"stage\" IN ('PENDING', 'FAILED') "

	services.DBCPsql.Raw(query).Find(&disbursements)
	ctx.JSON(iris.StatusOK, iris.Map{"data": disbursements})
}
