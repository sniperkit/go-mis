package role

import "time"

type Role struct {
	ID          uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name        string     `gorm:"column:name" json:"name"`
	Description string     `gorm:"column:description" json:"description"`
	Config      string     `gorm:"column:config" json:"config" sql:"json"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
