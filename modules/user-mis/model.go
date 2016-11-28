package userMis

import "time"

type UserMis struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Fullname  string     `gorm:"column:fullname" json:"fullname"`
	Username  string     `gorm:"column:_username" json:"username"`
	Password  string     `gorm:"column:_password" json:"password"`
	PicUrl    string     `gorm:"column:picUrl" json:"picUrl"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
