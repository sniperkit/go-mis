package loanMonitoring

import "time"

type LoanMonitoring struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Version   string     `gorm:"column:version" json:"version"`
	Raw       string     `gorm:"column:raw" json:"raw" sql:"json"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
