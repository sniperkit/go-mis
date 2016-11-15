package account

import "time"

type Account struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type      string     `gorm:"column:type" json:"type"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
