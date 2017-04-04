package accountTransactionCredit

import (
	"time"

	"bitbucket.org/go-mis/services"
)

type AccountTransactionCredit struct {
	ID              uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type            string     `gorm:"column:type" json:"type"` // CASHOUT, INVEST
	TransactionDate time.Time  `gorm:"column:transactionDate" json:"transactionDate"`
	Amount          float64    `gorm:"column:amount" json:"amount"`
	Remark          string     `gorm:"column:remark" json:"remark"`
	CreatedAt       time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type Total struct {
	Amount float64 `gorm:"column:amount"`
}

// GetTotalAccountTransactionCredit - sum account_transaction_credit based on accountID
func GetTotalAccountTransactionCredit(accountID uint64) float64 {
	query := "SELECT SUM(account_transaction_credit.amount) AS \"amount\" FROM account "
	query += "JOIN r_account_transaction_credit ON r_account_transaction_credit.\"accountId\" = account.id "
	query += "JOIN account_transaction_credit ON account_transaction_credit.id = r_account_transaction_credit.\"accountTransactionCreditId\" "
	query += "WHERE account.id = ? AND account_transaction_credit.\"deletedAt\" IS NULL "

	totalSchema := Total{}
	services.DBCPsql.Raw(query, accountID).Scan(&totalSchema)

	return totalSchema.Amount
}
