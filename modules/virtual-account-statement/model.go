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
