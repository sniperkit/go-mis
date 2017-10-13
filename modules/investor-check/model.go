package investorCheck

import "time"

type (
	InvestorCheck struct {
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
		ActivationDate         *time.Time `gorm:"column:activationDate" json:"activationDate"`
		DeclinedDate           *time.Time `gorm:"column:declinedDate" json:"declinedDate"`
		IsActivated            *bool      `gorm:"column:isActivated" json:"isActivated"`
		IsVerified             *bool      `gorm:"column:isVerified" json:"isVerified"`
		IsDeclined             *bool      `gorm:"column:isDeclined" json:"isDeclined"`
		Status                 string     `json:"status"`
		Username               string     `gorm:"column:username" json:"email"`
		RowsFullCount          int        `gorm:"column:full_count"`
	}

	InvestorNumber struct {
		ID         uint64 `gorm:"primary_key" gorm:"column:id" json:"id"`
		InvestorNo int    `gorm:"column:investorNo" json:"investorNo"`
	}

	DataTable struct {
		Columns   []DataColumn `json:"columns"`
		OrderInfo []OrderInfo  `json:"order"`
		Start     int          `json:"start"`
		Length    int          `json:"length"`
		Search    Search       `json:"search"`
		Draw      int          `json:"draw"`
	}

	DataColumn struct {
		Data       string `json:"data"`
		Name       string `json:"name"`
		Searchable bool   `json:"searchable"`
		Orderable  bool   `json:"orderable"`
		Search     `json:"search"`
	}

	Search struct {
		Value string `json:"value"`
		Regex bool   `json:"regex"`
	}

	OrderInfo struct {
		Column uint64 `json:"column"`
		Dir    string `json:"dir"`
	}
)
