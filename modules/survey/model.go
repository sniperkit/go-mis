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

type Survey struct {
	ID               uint64     `gorm:"column:id" json:"id"`
	BranchID         uint64     `gorm:"column:branchId" json:"branchId"`
	GroupID          uint64     `gorm:"column:groupId" json:"groupId"`
	AgentID          uint64     `gorm:"column:agentId" json:"agentId"`
	UUID             string     `gorm:"column:uuid" json:"uuid"`
	Fullname         string     `gorm:"column:fullname" json:"fullname"`
	CreditScoreGrade string     `gorm:"column:creditScoreGrade" json:"creditScoreGrade"`
	CreditScoreValue float64    `gorm:"column:creditScoreValue" json:"creditScoreValue"`
	Raw              string     `gorm:"column:_raw" json:"_raw" sql:"type:JSONB NULL DEFAULT '{}'::JSONB"`
	IsMigrate        bool       `gorm:"column:isMigrate" json:"isMigrate" sql:"DEFAULT:false"`
	IsApprove        bool       `gorm:"column:isApprove" json:"isApprove" sql:"DEFAULT:false"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
