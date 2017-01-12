package loanHistory

import "time"

type LoanHistory struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	StageFrom string     `gorm:"column:stageFrom" json:"stageFrom"`
	StageTo   string     `gorm:"column:stageTo" json:"stageTo"`
	Remark    string     `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
