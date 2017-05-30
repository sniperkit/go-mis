package campaign

import "time"

//Campaign - campaign database structure
type Campaign struct {
	ID           uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Amount       uint64     `gorm:"column:amount" json:"amount"`
	Name         string     `gorm:"column:name" json:"name"`
	Type         string     `gorm:"column:type" json:"type"`
	Description  string     `gorm:"column:description" json:"description"`
	MinimalValue int64      `gorm:"column:minimalValue" json:"minimalValue"`
	MaximalValue int64      `gorm:"column:maximalValue" json:"maximalValue"`
	Conversion   string     `gorm:"column:conversion" json:"conversion"`
	IsActive     bool       `gorm:"column:isActive" json:"isActive"`
	IsDebit      bool       `gorm:"column:isDebit" json:"isDebit"`
	StartDate    time.Time  `gorm:"column:startDate" json:"startDate"`
	EndDate      time.Time  `gorm:"column:endDate" json:"endDate"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
