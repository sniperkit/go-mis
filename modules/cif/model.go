package cif

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"bitbucket.org/go-mis/modules/investor"
)

type UpdateInvestor struct {
	Cif      Cif               `json:"cif"`
	Investor investor.Investor `json:"investor"`
}

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
	IdCardValidDate     string     `gorm:"column:idCardValidDate" json:"idCardValidDate"`
	IdCardFilename      string     `gorm:"column:idCardFilename" json:"idCardFilename"`
	TaxCardNo           string     `gorm:"column:taxCardNo" json:"taxCardNo"`
	TaxCardFilename     string     `gorm:"column:taxCardFilename" json:"taxCardFilename"`
	MaritalStatus       string     `gorm:"column:maritalStatus" json:"maritalStatus"`
	MotherName          string     `gorm:"column:mothersName" json:"mothersName"`
	Religion            string     `gorm:"column:religion" json:"religion"`
	Address             string     `gorm:"column:address" json:"address"`
	RT                  string     `gorm:"column:rt" json:"rt"`
	RW                  string     `gorm:"column:rw" json:"rw"`
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
	IsActivated         *bool      `gorm:"column:isActivated" json:"isActivated"`
	IsValidated         *bool      `gorm:"column:isValidated" json:"isValidated"`
	IsVerified          *bool      `gorm:"column:isVerified" json:"isVerified"`
	IsDeclined          *bool      `gorm:"column:isDeclined" json:"isDeclined"`
	CreatedAt           time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt           *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type CifBorrower struct {
	BorrowerId          string  `gorm:"column:borrowerId" json:"borrowerId"`
	BorrowerNo          string  `gorm:"column:borrowerNo" json:"borrowerNo"`
	IsCheckedTerm       *bool   `gorm:"column:isCheckedTerm" json:"isCheckedTerm"`
	IsCheckedPrivacy    *bool   `gorm:"column:IsCheckedPrivacy" json:"isCheckedPrivacy"`
	Village             string  `gorm:"column:village" json:"village"`
	CifNumber           uint64  `gorm:"column:cifNumber" json:"cifNumber"`
	Username            string  `gorm:"column:username" json:"username"`
	Password            string  `gorm:"column:password" json:"password"`
	Name                string  `gorm:"column:name" json:"name"`
	Gender              string  `gorm:"column:gender" json:"gender"`
	PlaceOfBirth        string  `gorm:"column:placeOfBirth" json:"placeOfBirth"`
	DateOfBirth         string  `gorm:"column:dateOfBirth" json:"dateOfBirth"`
	IdCardNo            string  `gorm:"column:idCardNo" json:"idCardNo"`
	IdCardValidDate     string  `gorm:"column:idCardValidDate" json:"idCardValidDate"`
	IdCardFilename      string  `gorm:"column:idCardFilename" json:"idCardFilename"`
	TaxCardNo           string  `gorm:"column:taxCardNo" json:"taxCardNo"`
	TaxCardFilename     string  `gorm:"column:taxCardFilename" json:"taxCardFilename"`
	MaritalStatus       string  `gorm:"column:maritalStatus" json:"maritalStatus"`
	MotherName          string  `gorm:"column:mothersName" json:"mothersName"`
	Religion            string  `gorm:"column:religion" json:"religion"`
	Address             string  `gorm:"column:address" json:"address"`
	RT                  string  `gorm:"column:rt" json:"rt"`
	RW                  string  `gorm:"column:rw" json:"rw"`
	Kelurahan           string  `gorm:"column:kelurahan" json:"kelurahan"`
	Kecamatan           string  `gorm:"column:kecamatan" json:"kecamatan"`
	City                string  `gorm:"column:city" json:"city"`
	Province            string  `gorm:"column:province" json:"province"`
	Nationality         string  `gorm:"column:nationality" json:"nationality"`
	Zipcode             string  `gorm:"column:zipcode" json:"zipcode"`
	PhoneNo             string  `gorm:"column:phoneNo" json:"phoneNo"`
	CompanyName         string  `gorm:"column:companyName" json:"companyName"`
	CompanyAddress      string  `gorm:"column:companyAddress" json:"companyAddress"`
	Occupation          string  `gorm:"column:occupation" json:"occupation"`
	Income              float64 `gorm:"column:income" json:"income"`
	IncomeSourceFund    string  `gorm:"column:incomeSourceFund" json:"incomeSourceFund"`
	IncomeSourceCountry string  `gorm:"column:incomeSourceCountry" json:"incomeSourceCountry"`
	IsActivated         *bool   `gorm:"column:isActivated" json:"isActivated"`
	IsValidated         *bool   `gorm:"column:isValidated" json:"isValidated"`
	IsVerified          *bool   `gorm:"column:isVerified" json:"isVerified"`
}

type CifInvestor struct {
	InvestorID          string  `gorm:"column:investorId" json:"investorId"`
	IsCheckedTerm       *bool   `gorm:"column:isCheckedTerm" json:"isCheckedTerm"`
	IsCheckedPrivacy    *bool   `gorm:"column:isCheckedPrivacy" json:"isCheckedPrivacy"`
	InvestorNo          string  `gorm:"column:investorNo" json:"investorNo"`
	IsInstitutional     *bool   `gorm:"column:isInstitutional" json:"isInstitutional"`
	BankName            string  `gorm:"column:bankName" json:"bankName"`
	BankBranch          string  `gorm:"column:bankBranch" json:"bankBranch"`
	BankAccountName     string  `gorm:"column:bankAccountName" json:"bankAccountName"`
	BankAccountNo       string  `gorm:"column:bankAccountNo" json:"bankAccountNo"`
	CifID               uint64  `gorm:"column:cifId" json:"cifId"`
	CifNumber           uint64  `gorm:"column:cifNumber" json:"cifNumber"`
	Username            string  `gorm:"column:username" json:"username"`
	Password            string  `gorm:"column:password" json:"password"`
	Name                string  `gorm:"column:name" json:"name"`
	Gender              string  `gorm:"column:gender" json:"gender"`
	PlaceOfBirth        string  `gorm:"column:placeOfBirth" json:"placeOfBirth"`
	DateOfBirth         string  `gorm:"column:dateOfBirth" json:"dateOfBirth"`
	IdCardNo            string  `gorm:"column:idCardNo" json:"idCardNo"`
	IdCardValidDate     string  `gorm:"column:idCardValidDate" json:"idCardValidDate"`
	IdCardFilename      string  `gorm:"column:idCardFilename" json:"idCardFilename"`
	TaxCardNo           string  `gorm:"column:taxCardNo" json:"taxCardNo"`
	TaxCardFilename     string  `gorm:"column:taxCardFilename" json:"taxCardFilename"`
	MaritalStatus       string  `gorm:"column:maritalStatus" json:"maritalStatus"`
	MotherName          string  `gorm:"column:mothersName" json:"mothersName"`
	Religion            string  `gorm:"column:religion" json:"religion"`
	Address             string  `gorm:"column:address" json:"address"`
	RT                  string  `gorm:"column:rt" json:"rt"`
	RW                  string  `gorm:"column:rw" json:"rw"`
	Kelurahan           string  `gorm:"column:kelurahan" json:"kelurahan"`
	Kecamatan           string  `gorm:"column:kecamatan" json:"kecamatan"`
	City                string  `gorm:"column:city" json:"city"`
	Province            string  `gorm:"column:province" json:"province"`
	Nationality         string  `gorm:"column:nationality" json:"nationality"`
	Zipcode             string  `gorm:"column:zipcode" json:"zipcode"`
	PhoneNo             string  `gorm:"column:phoneNo" json:"phoneNo"`
	CompanyName         string  `gorm:"column:companyName" json:"companyName"`
	CompanyAddress      string  `gorm:"column:companyAddress" json:"companyAddress"`
	Occupation          string  `gorm:"column:occupation" json:"occupation"`
	Income              float64 `gorm:"column:income" json:"income"`
	IncomeSourceFund    string  `gorm:"column:incomeSourceFund" json:"incomeSourceFund"`
	IncomeSourceCountry string  `gorm:"column:incomeSourceCountry" json:"incomeSourceCountry"`
	IsActivated         *bool   `gorm:"column:isActivated" json:"isActivated"`
	IsValidated         *bool   `gorm:"column:isValidated" json:"isValidated"`
	IsVerified          *bool   `gorm:"column:isVerified" json:"isVerified"`
}

type CifInvestorBorrower struct {
	ID          uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CifNumber   uint64 `gorm:"column:cifNumber" json:"cifNumber"`
	Name        string `gorm:"column:name" json:"name"`
	IsActivated *bool  `gorm:"column:isActivated" json:"isActivated"`
	IsValidated *bool  `gorm:"column:isValidated" json:"isValidated"`
	InvestorID  uint64 `gorm:"column:investorId" json:"investorId"`
	BorrowerID  uint64 `gorm:"column:borrowerId" json:"borrowerId"`
	IsBorrower  *bool  `gorm:"column:isBorrower" json:"isBorrower"`
	IsInvestor  *bool  `gorm:"column:isInvestor" json:"isInvestor"`
}

type CifSummary struct {
	TotalRegisteredCif uint64 `gorm:"column:totalRegisteredCif" json:"totalRegisteredCif"`
	TotalInvestor      uint64 `gorm:"column:totalInvestor" json:"totalInvestor"`
	TotalBorrower      uint64 `gorm:"column:totalBorrower" json:"totalBorrower"`
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
