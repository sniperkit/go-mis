package userMis

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type UserMis struct {
	ID          uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Fullname    string     `gorm:"column:fullname" json:"fullname"`
	Username    string     `gorm:"column:_username" json:"username"`
	Password    string     `gorm:"column:_password" json:"password"`
	PhoneNo     string     `gorm:"column:phoneNo" json:"phoneNo"`
	PicUrl      string     `gorm:"column:picUrl" json:"picUrl"`
	IsSuspended *bool      `gorm:"column:isSuspended" json:"isSuspended"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

func (u *UserMis) BeforeCreate() (err error) {
	if u.Password != "" {
		bytePassword := []byte(u.Password)
		sha256Bytes := sha256.Sum256(bytePassword)
		u.Password = hex.EncodeToString(sha256Bytes[:])
	}

	return
}

func (u *UserMis) BeforeUpdate() (err error) {
	if u.Password != "" {
		bytePassword := []byte(u.Password)
		sha256Bytes := sha256.Sum256(bytePassword)
		u.Password = hex.EncodeToString(sha256Bytes[:])
	}

	return
}

type UserMisAreaBranchRole struct {
	UserMisId   uint64 `gorm:"column:userMisId" json:"userMisId"`
	Fullname    string `gorm:"column:fullname" json:"fullname"`
	Username    string `gorm:"column:_username" json:"username"`
	PicUrl      string `gorm:"column:picUrl" json:"picUrl"`
	Password    string `gorm:"column:_password" json:"password"`
	PhoneNo    	string `gorm:"column:phoneNo" json:"phoneNo"`
	IsSuspended *bool  `gorm:"column:isSuspended" json:"isSuspended"`
	Area        string `gorm:"column:area" json:"area"`
	Branch      string `gorm:"column:branch" json:"branch"`
	Role        string `gorm:"column:role" json:"role"`
}
