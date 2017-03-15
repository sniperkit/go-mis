package account

import (
	"time"

	accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Account{})
	services.BaseCrudInit(Account{}, []Account{})
}

type ParamTopup struct {
	AccountID uint64  `json:"accountId"`
	Amount    float64 `json:"amount"`
}

type ParamBalance struct {
	Credit float64 `gorm:"column:credit"`
	Debit  float64 `gorm:"column:debit"`
}

// DoTopup - do topup
func DoTopup(ctx *iris.Context) {
	paramTopup := ParamTopup{}

	if err := ctx.ReadJSON(&paramTopup); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	userMis := ctx.Get("USER_MIS").(userMis.UserMis)
	remark := "TOPUP MANUAL by " + userMis.Fullname
	accountTransactionDebitSchema := &accountTransactionDebit.AccountTransactionDebit{Type: "TOPUP", TransactionDate: time.Now(), Amount: paramTopup.Amount, Remark: remark}
	services.DBCPsql.Create(accountTransactionDebitSchema)

	rAccountTransactionDebitSchema := &r.RAccountTransactionDebit{AccountId: paramTopup.AccountID, AccountTransactionDebitId: accountTransactionDebitSchema.ID}
	services.DBCPsql.Create(&rAccountTransactionDebitSchema)

	// query := "SELECT coalesce(sum(account_transaction_credit.amount), 0) AS \"credit\", coalesce(sum(account_transaction_debit.amount), 0) AS \"debit\" "
	// query += "FROM account "
	// query += "LEFT JOIN r_account_transaction_credit ON r_account_transaction_credit.\"accountId\" = account.id "
	// query += "LEFT JOIN account_transaction_credit ON account_transaction_credit.id = r_account_transaction_credit.\"accountTransactionCreditId\" "
	// query += "LEFT JOIN r_account_transaction_debit ON r_account_transaction_debit.\"accountId\" = account.id "
	// query += "LEFT JOIN account_transaction_debit ON account_transaction_debit.id = r_account_transaction_debit.\"accountTransactionDebitId\" "
	// query += "WHERE account.id = ? "
	// query += "AND account_transaction_debit.\"deletedAt\" IS NULL AND account_transaction_credit.\"deletedAt\" IS NULL "

	// paramBalance := ParamBalance{}
	// services.DBCPsql.Raw(query, paramTopup.AccountID).Scan(&paramBalance)

	// totalBalance := paramBalance.Debit - paramBalance.Credit

	// services.DBCPsql.Table("account").Where("id = ?", paramTopup.AccountID).Updates(Account{TotalDebit: paramBalance.Debit, TotalCredit: paramBalance.Credit, TotalBalance: totalBalance})

	totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebit(paramTopup.AccountID)
	totalCredit := accountTransactionCredit.GetTotalAccountTransactionCredit(paramTopup.AccountID)

	totalBalance := totalDebit - totalCredit

	services.DBCPsql.Table("account").Where("id = ?", paramTopup.AccountID).Updates(Account{TotalDebit: totalDebit, TotalCredit: totalCredit, TotalBalance: totalBalance})

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}
