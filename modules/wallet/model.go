package wallet

import "time"

type Wallet struct {
	ID           uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	TotalDebit   float64    `gorm:"column:totalDebit" json:"totalDebit"`
	TotalCredit  float64    `gorm:"column:totalCredit" json:"totalCredit"`
	TotalBalance float64    `gorm:"column:totalBalance" json:"totalBalance"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
