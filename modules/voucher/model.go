package voucher

import (
	"time"
)

type Voucher struct {
	ID          uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Amount      float64    `gorm:"column:amount" json:"amount"`
	VoucherNo   string     `gorm:"column:voucherNo" json:"voucherNo"`
	Description string     `gorm:"column:description" json:"description"`
	StarDate    *time.Time `gorm:"column:startDate" json:"startDate"`
	EndDate     *time.Time `gorm:"column:endDate" json:"endDate"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
	IsPersonal  *bool      `gorm:"column:isPersonal" json:"isPersonal"`
}
