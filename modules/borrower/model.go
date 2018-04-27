package borrower

import "time"

type Borrower struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	IsCheckedTerm    *bool      `gorm:"column:isCheckedTerm" json:"isCheckedTerm"`
	IsCheckedPrivacy *bool      `gorm:"column:isCheckedPrivacy" json:"isCheckedPrivacy"`
	BorrowerNo       string     `gorm:"column:borrowerNo" json:"borrowerNo"`
	Village          string     `gorm:"column:village" json:"village"`
	Education        string     `gorm:"column:education" json:"education"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
	DODate           *time.Time `gorm:"column:doDate" json:"doDate"`
}

type ProspectiveAvaraBorrower struct {
	BorrowerID      uint64  `gorm:"column:borrowerId" json:"borrowerId"`
	Name            string  `gorm:"column:name" json:"name"`
	GroupName       string  `gorm:"column:groupName" json:"groupName"`
	BranchID        uint64  `gorm:"column:branchId" json:"branchId"`
	Week            int     `gorm:"column:week" json:"week"`
	TotalPar        int     `gorm:"column:totalPar" json:"totalPar"`
	TotalTR         int     `gorm:"column:totalTR" json:"totalTR"`
	TotalPresence   int     `gorm:"column:totalPresence" json:"totalPresence"`
	TotalAvara      int     `gorm:"column:totalAvara" json:"totalAvara"`
	PresenceRatio   float64 `gorm:"column:presenceRatio" json:"presenceRatio"`
	Status          string  `gorm:"column:status" json:"status"`
}

type BorrowerDetail struct {
	ID           uint64 `gorm:"column:id" json:"id"`
	BorrowerName string `gorm:"column:borrowerName" json:"borrowerName"`
	PlaceOfBirth string `gorm:"column:placeOfBirth" json:"placeOfBirth"`
	DateOfBirth  string `gorm:"column:dateOfBirth" json:"dateOfBirth"`
	IDCardNo     string `gorm:"column:idCardNo" json:"idCardNo"`
	IDCardFilename string `gorm:"column:idCardFilename" json:"idCardFilename"`
	MaritalStatus string `gorm:"column:maritalStatus" json:"maritalStatus"`
	MotherName   string `gorm:"column:mothersName" json:"motherName"`
	Religion     string `gorm:"column:religion" json:"religion"`
	Address      string `gorm:"column:address" json:"address"`
	Kelurahan    string `gorm:"column:kelurahan" json:"kelurahan"`
	Kecamatan    string `gorm:"column:kecamatan" json:"kecamatan"`
	Occupation   string `gorm:"column:occupation" json:"occupation"`
	Income       uint64 `gorm:"column:income" json:"income"`
	RT           string `gorm:"column:rt" json:"rt"`
	RW           string `gorm:"column:rw" json:"rw"`
	BranchName   string `gorm:"column:branchName" json:"branchName"`
	GroupName    string `gorm:"column:groupName" json:"groupName"`
	PhoneNo      string `gorm:"column:phoneNo" json:"phoneNo"`
	FoName       string `gorm:"column:foName" json:"foName"`
	AgentName    string `gorm:"column:agentName" json:"agentName"`
	AgentAddress string `gorm:"column:agentAddress" json:"agentAddress"`
	AgentPhone	 string `gorm:"column:agentPhone" json:"agentPhone"`
}

type BorrowerAvaraRequest struct {
	BorrowerID  []uint64  `json:"borrowerId"`
	AgentID     uint64  `json:"agentId"`
}

var QUEUE_CREATE_AVARA_SURVEY = "QUEUE_CREATE_AVARA_SURVEY"