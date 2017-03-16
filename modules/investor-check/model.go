package investorCheck

import "time"

type InvestorCheck struct {
	ID              uint64     `gorm:"primary_key" gorm:"column:id" json:"_id"`
	Name            string     `gorm:"column:name" json:"name"`
	IDCardNo        string     `gorm:"column:idCardNo" json:"idCardNo"`
	IDCardFilename  string     `gorm:"column:idCardFilename" json:"idCardFilename"`
	TaxCardNo       string     `gorm:"column:taxCardNo" json:"taxCardNo"`
	TaxCardFilename string     `gorm:"column:taxCardFilename" json:"taxCardFilename"`
	CreatedAt       time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
