package investorCheck

import "time"

type InvestorCheck struct {
	ID                     uint64     `gorm:"primary_key" gorm:"column:id" json:"_id"`
	Name                   string     `gorm:"column:name" json:"name"`
	PhoneNo                string     `gorm:"column:phoneNo" json:"phoneNo"`
	IDCardNo               string     `gorm:"column:idCardNo" json:"idCardNo"`
	IDCardFilename         string     `gorm:"column:idCardFilename" json:"idCardFilename"`
	TaxCardNo              string     `gorm:"column:taxCardNo" json:"taxCardNo"`
	TaxCardFilename        string     `gorm:"column:taxCardFilename" json:"taxCardFilename"`
	BankAccountName        string     `gorm:"column:bankAccountName" json:"bankAccountName"`
	IsValidated            *bool      `gorm:"column:isValidated" json:"isValidated"`
	VirtualAccountBankName string     `gorm:"column:virtualAccountBankName" json:"virtualAccountBankName"`
	VirtualAccountNumber   string     `gorm:"column:virtualAccountNumber" json:"virtualAccountNumber"`
	InvestorNo             uint64     `gorm:"column:investorNo" json:"investorNo"`
	CreatedAt              time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt              time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt              *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
	IsActivated            *bool      `gorm:"column:isActivated" json:"isActivated"`
	IsVerified             *bool      `gorm:"column:isVerified" json:"isVerified"`
	IsDeclined             *bool      `gorm:"column:isDeclined" json:"isDeclined"`
	Status                 string     `json:"status"`
}

type InvestorNumber struct {
	ID         uint64 `gorm:"primary_key" gorm:"column:id" json:"id"`
	InvestorNo int    `gorm:"column:investorNo" json:"investorNo"`
}
