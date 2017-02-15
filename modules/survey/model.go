package survey

import "time"

type AFields struct {
	Key        string    `gorm:"column:key" json:"key"`
	Val        string    `gorm:"column:val" json:"val"`
	AnswerId   int64     `gorm:"column:answer_id" json:"answerId"`
	IsMigrated bool      `gorm:"is_migrated" json:"isMigrated"`
	IsApprove  bool      `gorm:"is_approve" json:"isApprove"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
