package investor

import "time"

type Investor struct {
	ID                       uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	IsInstitutional          *bool      `gorm:"column:isInstitutional" json:"isInstitutional"`
	IsCheckedTerm            *bool      `gorm:"column:isCheckedTerm" json:"isCheckedTerm"`
	IsCheckedPrivacy         *bool      `gorm:"column:isCheckedPrivacy" json:"isCheckedPrivacy"`
	InvestorNo               uint64     `gorm:"column:investorNo" json:"investorNo"`
	BankName                 string     `gorm:"column:bankName" json:"bankName"`
	BankBranch               string     `gorm:"column:bankBranch" json:"bankBranch"`
	BankAccountName          string     `gorm:"column:bankAccountName" json:"bankAccountName"`
	BankAccountNo            string     `gorm:"column:bankAccountNo" json:"bankAccountNo"`
	BorrowerCriteria         string     `gorm:"column:borrowerCriteria" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	IsBorrowerCriteriaActive *bool      `gorm:"column:isBorrowerCriteriaActive" json:"isBorrowerCriteriaActive"`
	CreatedAt                time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type InvestorWithoutVaSchema struct {
	Name       string `gorm:"column:name" json:"name"`
	InvestorID uint64 `gorm:"column:investorId" json:"investorId"`
	InvestorNo uint64 `gorm:"column:investorNo" json:"investorNo"`
}
