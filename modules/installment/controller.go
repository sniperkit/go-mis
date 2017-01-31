package installment

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Installment{})
	services.BaseCrudInit(Installment{}, []Installment{})
}

// FetchAll - fetchAll installment data
func FetchAll(ctx *iris.Context) {
	installments := []InstallmentFetch{}

	query := "SELECT installment.\"id\", installment.\"paidInstallment\" "
	query += ", \"group\".\"name\" as \"groupName\", branch.\"name\" as \"branchName\"  "
	query += "FROM installment "
	query += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
	query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = r_loan_installment.\"loanId\" "
	query += "JOIN branch ON r_loan_branch.\"branchId\" = branch.\"id\" "
	query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = r_loan_installment.\"loanId\" "
	query += "JOIN \"group\" ON r_loan_group.\"groupId\" = \"group\".\"id\" "

	services.DBCPsql.Raw(query).Find(&installments)
	ctx.JSON(iris.StatusOK, iris.Map{"data": installments})
}
