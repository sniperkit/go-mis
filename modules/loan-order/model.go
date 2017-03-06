package loanOrder

import "time"

type LoanOrder struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	OrderNo   string     `gorm:"column:orderNo" json:"orderNo"`
	Remark    string     `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
