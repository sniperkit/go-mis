package virtualAccount

import "time"

type VirtualAccount struct {
	ID                 uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	BankName           string     `gorm:"column:bankName" json:"bankName"`
	VirtualAccountCode string     `gorm:"column:virtualAccountCode" json:"virtualAccountCode"`
	VirtualAccountNo   string     `gorm:"column:virtualAccountNo" json:"virtualAccountNo"`
	VirtualAccountName string     `gorm:"column:virtualAccountName" json:"virtualAccountName"`
	CreatedAt          time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt          *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
