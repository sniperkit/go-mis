package accountTransactionDebit

import "time"

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
