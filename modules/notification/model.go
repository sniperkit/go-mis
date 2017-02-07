package notification

import "time"

type Notification struct {
	ID          uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type        string     `gorm:"column:type" json:"type"`
	Message     string     `gorm:"column:message" json:"message"`
	IsRead      *bool      `gorm:"column:isRead" json:"isRead"`
	RedirectUrl string     `gorm:"column:redirectUrl" json:"redirectUrl"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type NotificationInput struct {
	SentTo      string `json:"sentTo"`
	Message     string `gorm:"column:message" json:"message"`
	RedirectUrl string `gorm:"column:redirectUrl" json:"redirectUrl"`
}
