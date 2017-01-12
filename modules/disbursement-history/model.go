package disbursementHistory

import "time"

type DisbursementHistory struct {
	ID                   uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	StageFrom            string     `gorm:"column:stageFrom" json:"stageFrom"`
	StageTo              string     `gorm:"column:stageTo" json:"stageTo"`
	Remark               string     `gorm:"column:remark" json:"remark"`
	LastDisbursementDate time.Time  `gorm:"column:lastDisbursementDate" json:"lastDisbursementDate"`
	NextDisbursementDate time.Time  `gorm:"column:nextDisbursementDate" json:"nextDisbursementDate"`
	CreatedAt            time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt            time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt            *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
