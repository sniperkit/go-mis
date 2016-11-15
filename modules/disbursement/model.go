package disbursement

import "time"

type Disbursement struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Status    string     `gorm:"column:status" json:"status"`
	Message   string     `gorm:"column:message" json:"message"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
