package disbursement

import "time"

type Disbursement struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	DisbursementDate time.Time  `gorm:"column:disbursementDate" json:"disbursementDate"`
	Stage            string     `gorm:"column:stage" json:"stage"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type DisbursementFetch struct {
	ID                uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Borrower          string    `gorm:"column:borrower" json:"borrower"`
	Group             string    `gorm:"column:group" json:"group"`
	Branch            string    `gorm:"column:branch" json:"branch"`
	SubmittedLoanDate time.Time `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
	DisbursementDate  time.Time `gorm:"column:disbursementDate" json:"disbursementDate"`
	Plafond           float64   `gorm:"column:plafond" json:"plafond"`
	Stage             string    `gorm:"column:stage" json:"stage"`
	CreatedAt         time.Time `gorm:"column:createdAt" json:"createdAt"`
}
