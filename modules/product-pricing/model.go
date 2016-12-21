package productPricing

import "time"

type ProductPricing struct {
	ID                 uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	ReturnOfInvestment float64    `gorm:"column:returnOfInvestment" json:"returnOfInvestment"`
	AdminitrationFee   float64    `gorm:"column:administrationFee" json:"administrationFee"`
	CreatedAt          time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt          *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
