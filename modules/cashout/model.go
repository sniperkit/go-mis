package cashout

import "time"

type Cashout struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CashoutID string     `gorm:"column:cashoutId" json:"cashoutId"`
	Amount    float64    `gorm:"column:amount" json:"amount"`
	Stage     string     `gorm:"column:stage" json:"stage"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
