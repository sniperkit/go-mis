package cif

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type Cif struct {
	ID                  uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CifNumber           uint64     `gorm:"column:cifNumber" json:"cifNumber"`
	Username            string     `gorm:"column:username" json:"username"`
	Password            string     `gorm:"column:password" json:"password"`
	Name                string     `gorm:"column:name" json:"name"`
	Gender              string     `gorm:"column:gender" json:"gender"`
	PlaceOfBirth        string     `gorm:"column:placeOfBirth" json:"placeOfBirth"`
	DateOfBirth         string     `gorm:"column:dateOfBirth" json:"dateOfBirth"`
	IdCardNo            string     `gorm:"column:idCardNo" json:"idCardNo"`
	IdCardValidDate     time.Time  `gorm:"column:idCardValidDate" json:"idCardValidDate"`
	IdCardFilename      string     `gorm:"column:idCardFilename" json:"idCardFilename"`
	TaxCardNo           string     `gorm:"column:taxCardNo" json:"taxCardNo"`
	TaxCardFilename     string     `gorm:"column:taxCardFilename" json:"taxCardFilename"`
	MaritalStatus       string     `gorm:"column:maritalStatus" json:"maritalStatus"`
	MotherName          string     `gorm:"column:mothersName" json:"mothersName"`
	Religion            string     `gorm:"column:religion" json:"religion"`
	Address             string     `gorm:"column:address" json:"address"`
	Kelurahan           string     `gorm:"column:kelurahan" json:"kelurahan"`
	Kecamatan           string     `gorm:"column:kecamatan" json:"kecamatan"`
	City                string     `gorm:"column:city" json:"city"`
	Province            string     `gorm:"column:province" json:"province"`
	Nationality         string     `gorm:"column:nationality" json:"nationality"`
	Zipcode             string     `gorm:"column:zipcode" json:"zipcode"`
	PhoneNo             string     `gorm:"column:phoneNo" json:"phoneNo"`
	CompanyName         string     `gorm:"column:companyName" json:"companyName"`
	CompanyAddress      string     `gorm:"column:companyAddress" json:"companyAddress"`
	Occupation          string     `gorm:"column:occupation" json:"occupation"`
	Income              float64    `gorm:"column:income" json:"income"`
	IncomeSourceFund    string     `gorm:"column:incomeSourceFund" json:"incomeSourceFund"`
	IncomeSourceCountry string     `gorm:"column:incomeSourceCountry" json:"incomeSourceCountry"`
	IsActivated         bool       `gorm:"column:isActivated" json:"isActivated"`
	IsVAlidated         bool       `gorm:"column:isValidated" json:"isValidated"`
	CreatedAt           time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt           *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

func (c *Cif) BeforeCreate() (err error) {
	if c.Password != "" {
		bytePassword := []byte(c.Password)
		md5Bytes := md5.Sum(bytePassword)
		c.Password = hex.EncodeToString(md5Bytes[:])
	}

	return
}

func (c *Cif) BeforeUpdate() (err error) {
	if c.Password != "" {
		bytePassword := []byte(c.Password)
		md5Bytes := md5.Sum(bytePassword)
		c.Password = hex.EncodeToString(md5Bytes[:])
	}

	return
}
