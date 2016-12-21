package sector

import "time"

type Sector struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name      string     `gorm:"column:name" json:"name"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
