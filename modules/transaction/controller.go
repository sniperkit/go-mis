package transaction

import (
	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

// GetData - get data by transaction type
func GetData(ctx *iris.Context) {
	transactionType := ctx.Param("type")
	investorID := ctx.Param("investor_id")
	startDate := ctx.Param("start_date") + " 00:00:00"
	endDate := ctx.Param("end_date") + " 00:00:00"

	totalData := TotalData{}

	query := "SELECT account_transaction_debit.* FROM account_transaction_debit "
	query += "JOIN r_account_transaction_debit ON r_account_transaction_debit.\"accountTransactionDebitId\" = account_transaction_debit.id "
	query += "JOIN r_account_investor ON r_account_investor.\"accountId\" = r_account_transaction_debit.\"accountId\" "
	query += "WHERE r_account_investor.\"investorId\" = ? "
	query += "AND account_transaction_debit.\"type\" = ? "
	query += "AND account_transaction_debit.\"transactionDate\" BETWEEN ? AND ? "

	queryTotal := "SELECT count(account_transaction_debit.*) as \"totalRows\" FROM account_transaction_debit "
	queryTotal += "JOIN r_account_transaction_debit ON r_account_transaction_debit.\"accountTransactionDebitId\" = account_transaction_debit.id "
	queryTotal += "JOIN r_account_investor ON r_account_investor.\"accountId\" = r_account_transaction_debit.\"accountId\" "
	queryTotal += "WHERE r_account_investor.\"investorId\" = ? "
	queryTotal += "AND account_transaction_debit.\"type\" = ? "
	queryTotal += "AND account_transaction_debit.\"transactionDate\" BETWEEN ? AND ? "

	if ctx.URLParam("LIMIT") != "" {
		query += "LIMIT " + ctx.URLParam("LIMIT")
		queryTotal += "LIMIT " + ctx.URLParam("LIMIT")
	} else {
		query += "LIMIT 10 "
		queryTotal += "LIMIT 10 "
	}

	services.DBCPsql.Raw(queryTotal).Find(&totalData)

	atdSchema := []accountTransactionDebit.AccountTransactionDebit{}
	services.DBCPsql.Raw(query, investorID, transactionType, startDate, endDate).Scan(&atdSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      atdSchema,
	})
}
