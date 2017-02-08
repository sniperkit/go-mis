package borrower

import "time"

type Borrower struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	IsCheckedTerm    *bool      `gorm:"column:isCheckedTerm" json:"isCheckedTerm"`
	IsCheckedPrivacy *bool      `gorm:"column:isCheckedPrivacy" json:"isCheckedPrivacy"`
	BorrowerNo       string     `gorm:"column:borrowerNo" json:"borrowerNo"`
	Village          string     `gorm:"column:village" json:"village"`
	Education        string     `gorm:"column:education" json:"education"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
