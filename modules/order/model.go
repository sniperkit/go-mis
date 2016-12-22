package order

import "time"

type Order struct {
	ID          uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	OrderNo     string     `gorm:"column:orderNo" json:"orderNo"`
	GrossAmount float64    `gorm:"column:grossAmount" json:"grossAmount"`
	NettAmount  float64    `gorm:"column:nettAmount" json:"nettAmount"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
