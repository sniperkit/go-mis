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
}
