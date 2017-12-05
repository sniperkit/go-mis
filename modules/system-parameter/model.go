package systemParameter

import (
	"time"
)

type (
	SystemParameter struct {
		ID        uint64     `bson:"column:id"`
		Key       string     `bson:"key" json:"key"`
		Value     string     `bson:"value" json:"value"`
		CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
		UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
		DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
	}
	Swift struct {
		Code   string `gorm:"column:code" json:"code"`
		Bank   string `gorm:"column:bank" json:"bank"`
		City   string `gorm:"column:city" json:"city"`
		Branch string `gorm:"column:branch" json:"branch"`
		Prefix string `gorm:"column:prefix" json:"prefix"`
	}
)
