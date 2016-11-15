package walletTransaction

import "time"

type WalletTransaction struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Debit     float64    `gorm:"column:debit" json:"debit"`
	Credit    float64    `gorm:"column:credit" json:"credit"`
	Balance   float64    `gorm:"column:balance" json:"balance"`
	Remark    string     `gorm:"column:remark" json:"remark"`
	Status    string     `gorm:"column:status" json:"status"`
	Type      string     `gorm:"column:type" json:"type"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
