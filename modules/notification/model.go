package notification

import "time"

type Notification struct {
	ID          uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type        string     `gorm:"column:type" json:"type"`
	Message     string     `gorm:"column:message" json:"message"`
	IsRead      bool       `gorm:"column:isRead" json:"isRead"`
	RedirectUrl string     `gorm:"column:redirectUrl" json:"redirectUrl"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
