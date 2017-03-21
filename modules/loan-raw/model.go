package loanRaw

import "time"

type LoanRaw struct {
	ID        uint64      `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanID    uint64      `gorm:"column:loanId" json:"loanId"`
	Version   string      `gorm:"column:_version" json:"_version"`
	Raw       interface{} `gorm:"column:_raw" json:"_raw" sql:"type:JSONB NULL DEFAULT '{}'::JSONB"`
	CreatedAt time.Time   `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time   `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time  `gorm:"column:deletedAt" json:"deletedAt"`
}
