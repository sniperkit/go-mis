package installment

import (
	"strconv"

	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	installmentHistory "bitbucket.org/go-mis/modules/installment-history"
	"bitbucket.org/go-mis/modules/r"
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

// SubmitInstallment - submit installment data
func SubmitInstallment(ctx *iris.Context) {
	installment := Installment{}
	tempLoanID := ctx.Param("loan_id")

	loanID, err := strconv.ParseUint(tempLoanID, 10, 64)

	if err != nil {
		ctx.JSON(iris.StatusExpectationFailed, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = ctx.ReadJSON(&installment)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	services.DBCPsql.Table("installment").Create(&installment)

	installmentHistoryData := installmentHistory.InstallmentHistory{StageFrom: installment.Stage, StageTo: installment.Stage}
	services.DBCPsql.Table("r_installment_history").Create(&installmentHistoryData)

	rInstallmentHistoryData := r.RInstallmentHistory{InstallmentId: installment.ID, InstallmentHistoryId: installmentHistoryData.ID}
	services.DBCPsql.Table("r_installment_history").Create(&rInstallmentHistoryData)

	accountTransactionDebitData := accountTransactionDebit.AccountTransactionDebit{Type: "INSTALLMENT", Amount: installment.PaidInstallment}
	services.DBCPsql.Table("account_transaction_debit").Create(&accountTransactionDebitData)

	rLoanInvestorProductPricing := r.RLoanInvestorProductPricing{}
	services.DBCPsql.Table("r_loan_investor_product_pricing").Where("\"loanId\" = ?", loanID).First(&rLoanInvestorProductPricing)

	rAccountInvestor := r.RAccountInvestor{}
	services.DBCPsql.Table("r_account_investor").Where("\"investorId\" = ?", rLoanInvestorProductPricing.InvestorId).First(&rAccountInvestor)

	rAccountTransactionDebitData := r.RAccountTransactionDebit{AccountId: rAccountInvestor.AccountId, AccountTransactionDebitId: accountTransactionDebitData.ID}
	go services.DBCPsql.Table("r_account_transaction_debit").Create(&rAccountTransactionDebitData)

	loanInstallmentData := r.RLoanInstallment{LoanId: loanID, InstallmentId: installment.ID}
	go services.DBCPsql.Table("r_loan_installment").Create(&loanInstallmentData)

	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status": "success",
		"data":   installment,
	})

}
