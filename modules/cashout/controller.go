package cashout

import (
	"strconv"
	"time"

	"bitbucket.org/go-mis/modules/account"
	accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	cashoutHistory "bitbucket.org/go-mis/modules/cashout-history"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Cashout{})
	services.BaseCrudInit(Cashout{}, []Cashout{})
}

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

// FetchAll - fetch cashout data
func FetchAll(ctx *iris.Context) {
	cashoutInvestors := []CashoutInvestor{}
	// stage := ctx.URLParam("stage")

	// query := "SELECT cashout.*, "
	// query += "cif.\"name\" as \"investorName\", "
	// query += "account.\"totalBalance\" as \"totalBalance\" "
	// // query += "r_cif_borrower.\"borrowerId\" IS NOT NULL as \"isBorrower\", r_cif_investor.\"investorId\" IS NOT NULL as \"isInvestor\" "
	// query += "FROM cashout "
	// query += "LEFT JOIN r_investor_cashout ON r_investor_cashout.\"cashoutId\" = cashout.\"id\" "
	// query += "LEFT JOIN investor ON investor.\"id\" = r_investor_cashout.\"investorId\" "
	// query += "LEFT JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.\"id\" "
	// query += "LEFT JOIN cif ON r_cif_investor.\"cifId\" = cif.\"id\" "
	// query += "LEFT JOIN r_account_investor ON r_account_investor.\"investorId\" = investor.\"id\" "
	// query += "LEFT JOIN account ON r_account_investor.\"accountId\" = account.\"id\" "

	query := "SELECT r_cif_investor.\"cifId\", r_cif_investor.\"investorId\", r_account_transaction_credit_cashout.\"cashoutId\", r_account_investor.\"accountId\", cif.name AS \"investorName\", cashout.\"cashoutId\" AS \"cashoutNo\", cashout.amount, account.\"totalDebit\", account.\"totalCredit\", account.\"totalBalance\", account_transaction_credit.\"type\",  account_transaction_credit.\"transactionDate\", account_transaction_credit.remark, cashout.stage "
	query += "FROM cashout "
	query += "JOIN r_account_transaction_credit_cashout ON r_account_transaction_credit_cashout.\"cashoutId\" = cashout.id "
	query += "JOIN r_account_transaction_credit ON r_account_transaction_credit.\"accountTransactionCreditId\" = r_account_transaction_credit_cashout.\"accountTransactionCreditId\" "
	query += "JOIN account_transaction_credit ON account_transaction_credit.id = r_account_transaction_credit_cashout.\"accountTransactionCreditId\" "
	query += "JOIN account ON account.id = r_account_transaction_credit.\"accountId\" "
	query += "JOIN r_account_investor ON r_account_investor.\"accountId\" = account.id "
	query += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = r_account_investor.\"investorId\" "
	query += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "

	queryTotalData := "SELECT count(*) AS \"totalRows\" "
	queryTotalData += "FROM cashout "
	queryTotalData += "JOIN r_account_transaction_credit_cashout ON r_account_transaction_credit_cashout.\"cashoutId\" = cashout.id "
	queryTotalData += "JOIN r_account_transaction_credit ON r_account_transaction_credit.\"accountTransactionCreditId\" = r_account_transaction_credit_cashout.\"accountTransactionCreditId\" "
	queryTotalData += "JOIN account_transaction_credit ON account_transaction_credit.id = r_account_transaction_credit_cashout.\"accountTransactionCreditId\" "
	queryTotalData += "JOIN account ON account.id = r_account_transaction_credit.\"accountId\" "
	queryTotalData += "JOIN r_account_investor ON r_account_investor.\"accountId\" = account.id "
	queryTotalData += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = r_account_investor.\"investorId\" "
	queryTotalData += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "

	// if len(strings.TrimSpace(stage)) == 0 {
	// 	services.DBCPsql.Raw(query).Find(&cashoutInvestors)
	// } else {
	// 	query += "WHERE cashout.\"stage\" = ? "
	// 	services.DBCPsql.Raw(query, stage).Find(&cashoutInvestors)
	// }

	totalData := TotalData{}

	services.DBCPsql.Raw(query).Find(&cashoutInvestors)
	services.DBCPsql.Raw(queryTotalData).Find(&totalData)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      cashoutInvestors,
	})
}

// UpdateStage - update cashout stage
func UpdateStage(ctx *iris.Context) {
	cashoutInvestorSchema := CashoutInvestor{}

	if err := ctx.ReadJSON(&cashoutInvestorSchema); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Failed to read param",
		})
		return
	}

	cashoutID, _ := strconv.ParseUint(ctx.Param("cashout_id"), 10, 64)
	stage := ctx.Param("stage")

	cashoutHistoryObj := &cashoutHistory.CashoutHistory{StageFrom: "PENDING", StageTo: stage}
	services.DBCPsql.Create(cashoutHistoryObj)

	rCashoutHistoryObj := &r.RCashoutHistory{CashoutId: cashoutID, CashoutHistoryId: cashoutHistoryObj.ID}
	services.DBCPsql.Create(rCashoutHistoryObj)

	services.DBCPsql.Table("cashout").Where("cashout.\"id\" = ?", cashoutID).UpdateColumn("stage", stage)

	cashoutNo := strconv.FormatUint(cashoutInvestorSchema.CashoutNo, 10)

	if stage == "SEND-TO-MANDIRI" {
		// TODO: hit go-banking for Mandiri Corporate Payable
		
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   iris.Map{},
		})
	} else if stage == "SUCCESS" {
		services.DBCPsql.Table("account_transaction_credit").Where("remark = ?", cashoutInvestorSchema.Remark).Update("remark", "CASHOUT #"+cashoutNo+" SUCCESS")
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   iris.Map{},
		})
	} else if stage == "CANCEL" {
		services.DBCPsql.Table("account_transaction_credit").Where("remark = ?", cashoutInvestorSchema.Remark).Update("remark", "CASHOUT #"+cashoutNo+" CANCEL")

		accountTransactionDebitSchema := &accountTransactionDebit.AccountTransactionDebit{Type: "REFUND", Remark: "REFUND CASHOUT #" + cashoutNo, Amount: cashoutInvestorSchema.Amount, TransactionDate: time.Now()}
		services.DBCPsql.Table("account_transaction_debit").Create(accountTransactionDebitSchema)

		rAccountTransactionDebitSchema := &r.RAccountTransactionDebit{AccountId: cashoutInvestorSchema.AccountID, AccountTransactionDebitId: accountTransactionDebitSchema.ID}
		services.DBCPsql.Table("r_account_transaction_debit").Create(rAccountTransactionDebitSchema)

		totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebit(cashoutInvestorSchema.AccountID)
		totalCredit := accountTransactionCredit.GetTotalAccountTransactionCredit(cashoutInvestorSchema.AccountID)

		totalBalance := totalDebit - totalCredit

		services.DBCPsql.Table("account").Where("id = ?", cashoutInvestorSchema.AccountID).Updates(account.Account{TotalDebit: totalDebit, TotalCredit: totalCredit, TotalBalance: totalBalance})

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   iris.Map{},
		})
	} else {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": "Bad request.",
		})
		return
	}
}
