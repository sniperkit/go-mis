package campaign

import "time"

type Campaign struct {
	ID          uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name        string     `gorm:"column:name" json:"name"`
	Description string     `gorm:"column:description" json:"description"`
	StartDate   time.Time  `gorm:"column:startDate" json:"startDate"`
	EndDate     time.Time  `gorm:"column:endDate" json:"endDate"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
