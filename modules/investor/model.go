package investor

import "time"

type Investor struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	IsInstitutional  bool       `gorm:"column:isInstitutional" json:"isInstitutional"`
	IsCheckedTerm    bool       `gorm:"column:isCheckedTerm" json:"isCheckedTerm"`
	IsCheckedPrivacy bool       `gorm:"column:isCheckedPrivacy" json:"isCheckedPrivacy"`
	InvestorNo       uint64     `gorm:"column:investorNo" json:"investorNo"`
	BankName         string     `gorm:"column:bankName" json:"bankName"`
	BankBranch       string     `gorm:"column:bankBranch" json:"bankBranch"`
	BankAccountName  string     `gorm:"column:bankAccountName" json:"bankAccountName"`
	BankAccountNo    uint64     `gorm:"column:bankAccountNo" json:"bankAccountNo"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
