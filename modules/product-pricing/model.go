package productPricing

import "time"

type ProductPricing struct {
	ID                 uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	ReturnOfInvestment float64    `gorm:"column:returnOfInvestment" json:"returnOfInvestment"`
	AdminitrationFee   float64    `gorm:"column:administrationFee" json:"administrationFee"`
	ServiceFee         float64    `gorm:"column:serviceFee" json:"serviceFee"`
	StartDate          *time.Time `gorm:"column:startDate" json:"startDate"`
	EndDate            *time.Time `gorm:"column:endDate" json:"endDate"`
	IsInstutitional    *bool      `gorm:"column:isInstitutional" json:"isInstitutional"`
	CreatedAt          time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt          *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type InvestorSearch struct {
	ID                uint64  `gorm:"column:id" json:"id"`
	Investor					string  `gorm:"column:name" json:"name"`
}
