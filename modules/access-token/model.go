package accessToken

import "time"

type AccessToken struct {
	ID          uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type        string     `gorm:"column:type" json:"type"` // WEB, ANDROID, IOS
	AccessToken string     `gorm:"column:accessToken" json:"accessToken"`
	ExpiredAt   *time.Time `gorm:"column:expiredAt" json:"expiredAt"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
