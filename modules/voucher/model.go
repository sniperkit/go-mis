package voucher

import (
	"time"
	"fmt"
)

type Voucher struct {
	ID        				uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Amount    				uint64    `gorm:"column:amount" json:"amount"`
	VoucherNo 				string    `gorm:"column:voucherrNo" json:"voucherNo"`
	Description 			string    `gorm:"column:description" json:"description"`
	StarDate 					time.Time `gorm:"column:startDate" json:"startDate"`
	EndDate 					time.Time `gorm:"column:endDate" json:"endDate"`
	CreatedAt 				time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdateAt 					time.Time `gorm:"column:updateAt" json:"updateAt"`
	IsPersonal				*bool 		`gorm:"column:isPersonal" json:"isPersonal"`
}
