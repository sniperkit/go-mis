package virtualAccountStatement

import "time"

type VirtualAccountStatement struct {
	ID              uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	TransactionDate *time.Time `gorm:"column:transactionDate" json:"transactionDate"`
	Currency        string     `gorm:"column:currency" json:"currency"`
	Amount          float64    `gorm:"column:amount" json:"amount"`
	CreatedAt       time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type VirtualAccountStatementSubmission struct {
	TransactionDate    *time.Time `gorm:"column:transactionDate" json:"transactionDate"`
	Amount             float64    `gorm:"column:amount" json:"amount"`
	BankName           string     `gorm:"column:bankName" json:"bankName"`
	VirtualAccountNo   string     `gorm:"column:virtualAccountNo" json:"virtualAccountNo"`
	VirtualAccountName string     `gorm:"column:virtualAccountName" json:"virtualAccountName"`
}
