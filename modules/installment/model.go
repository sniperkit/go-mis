package installment

import "time"

type Installment struct {
	ID             uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type           string     `gorm:"column:type" json:"type"`
	Principal      float64    `gorm:"column:principal" json:"principal"`
	ProfitAmf      float64    `gorm:"column:profitAmf" json:"profitAmf"`
	ProfitInvestor float64    `gorm:"column:profitInvestor" json:"profitInvestor"`
	Reserve        float64    `gorm:"column:reserve" json:"reserve"`
	CreatedAt      time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
