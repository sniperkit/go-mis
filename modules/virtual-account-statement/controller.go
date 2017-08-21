package virtualAccountStatement

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&VirtualAccountStatement{})
	services.BaseCrudInit(VirtualAccountStatement{}, []VirtualAccountStatement{})
}

type TransactionByBank struct {
	Total int64 `gorm:"column:total" json:"total"`
}

// GetVAStatement - get VA statement
func GetVAStatement(ctx *iris.Context) {
	query := "SELECT "
	query += "virtual_account_statement.\"transactionDate\", virtual_account_statement.amount, "
	query += "virtual_account.\"bankName\", virtual_account.\"virtualAccountNo\", virtual_account.\"virtualAccountName\" "
	query += "FROM virtual_account "
	query += "INNER JOIN r_virtual_account_statement ON r_virtual_account_statement.\"virtualAccountId\" = virtual_account.id "
	query += "INNER JOIN virtual_account_statement ON virtual_account_statement.id = r_virtual_account_statement.\"virtualAccountStatementId\" "
	query += "WHERE virtual_account_statement.\"deletedAt\" IS NULL AND virtual_account.\"deletedAt\" IS NULL"

	vaStatmentSubmissionSchema := []VirtualAccountStatementSubmission{}
	services.DBCPsql.Raw(query).Scan(&vaStatmentSubmissionSchema)

	queryTotalBCA := "SELECT COUNT(*) AS \"total\" "
	queryTotalBCA += "FROM virtual_account "
	queryTotalBCA += "INNER JOIN r_virtual_account_statement ON r_virtual_account_statement.\"virtualAccountId\" = virtual_account.id "
	queryTotalBCA += "INNER JOIN virtual_account_statement ON virtual_account_statement.id = r_virtual_account_statement.\"virtualAccountStatementId\" "
	queryTotalBCA += "WHERE virtual_account.\"bankName\" = 'BCA' AND virtual_account_statement.\"deletedAt\" IS NULL AND virtual_account.\"deletedAt\" IS NULL"

	totalBCA := TransactionByBank{}
	services.DBCPsql.Raw(queryTotalBCA).Scan(&totalBCA)

	queryTotalBRI := "SELECT COUNT(*) AS \"total\" "
	queryTotalBRI += "FROM virtual_account "
	queryTotalBRI += "INNER JOIN r_virtual_account_statement ON r_virtual_account_statement.\"virtualAccountId\" = virtual_account.id "
	queryTotalBRI += "INNER JOIN virtual_account_statement ON virtual_account_statement.id = r_virtual_account_statement.\"virtualAccountStatementId\" "
	queryTotalBRI += "WHERE virtual_account.\"bankName\" = 'BRI' AND virtual_account_statement.\"deletedAt\" IS NULL AND virtual_account.\"deletedAt\" IS NULL"

	totalBRI := TransactionByBank{}
	services.DBCPsql.Raw(queryTotalBRI).Scan(&totalBRI)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"transactionByBank": iris.Map{
				"BCA": totalBCA.Total,
				"BRI": totalBRI.Total,
			},
			"vaStatementDetail": vaStatmentSubmissionSchema,
		},
	})
}
