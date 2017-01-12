package adjustment

import "time"

type Adjustment struct {
	ID             uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type           string     `gorm:"column:type" json:"type"` // INSTALLMENT, DISBURSEMENT, ACCOUNT-CREDIT, ACCOUNT-DEBIT
	AmountBefore   float64    `gorm:"column:amountBefore" json:"amountBefore"`
	AmountToAdjust float64    `gorm:"column:amountToAdjust" json:"amountToAdjust"`
	AmountAfter    float64    `gorm:"column:amountAfter" json:"amountAfter"`
	Remark         string     `gorm:"column:remark" json:"remark"`
	CreatedAt      time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
