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
	IsDeclined             *bool      `gorm:"column:isDeclined"`
	VirtualAccountBankName string     `gorm:"column:virtualAccountBankName" json:"virtualAccountBankName"`
	VirtualAccountNumber   string     `gorm:"column:virtualAccountNumber" json:"virtualAccountNumber"`
	InvestorNo             uint64     `gorm:"column:investorNo" json:"investorNo"`
	Email                  string     `gorm:"column:username" json:"email"`
	ActivationDate         time.Time  `gorm:"column:activationDate" json:"activationDate"`
	DeclinedDate           time.Time  `gorm:"column:declinedDate" json:"declinedDate"`
	Status                 string     `json:"status"`
	CreatedAt              time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt              time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt              *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type InvestorNumber struct {
	ID         uint64 `gorm:"primary_key" gorm:"column:id" json:"id"`
	InvestorNo int    `gorm:"column:investorNo" json:"investorNo"`
}
