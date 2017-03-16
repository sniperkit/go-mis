package accountTransactionDebit

import (
	"time"

	"bitbucket.org/go-mis/services"
)

type AccountTransactionDebit struct {
	ID              uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type            string     `gorm:"column:type" json:"type"` // REFUND, REFERRAL, INSTALLMENT, TOPUP via VA
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

// GetTotalAccountTransactionDebit - sum account_transaction_credit based on accountID
func GetTotalAccountTransactionDebit(accountID uint64) float64 {
	query := "SELECT SUM(account_transaction_debit.amount) AS \"amount\" FROM account "
	query += "JOIN r_account_transaction_debit ON r_account_transaction_debit.\"accountId\" = account.id "
	query += "JOIN account_transaction_debit ON account_transaction_debit.id = r_account_transaction_debit.\"accountTransactionDebitId\" "
	query += "WHERE account.id = ? "

	totalSchema := Total{}
	services.DBCPsql.Raw(query, accountID).Scan(&totalSchema)

	return totalSchema.Amount
}
