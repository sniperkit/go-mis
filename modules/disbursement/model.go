package disbursement

import "time"

type Disbursement struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	DisbursementDate string     `gorm:"column:disbursementDate" json:"disbursementDate"`
	Stage            string     `gorm:"column:stage" json:"stage"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type DisbursementFetch struct {
	ID                uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	GroupId           uint64    `gorm:"column:groupId" json:"groupId"`
	Group             string    `gorm:"column:group" json:"group"`
	BranchId          uint64    `gorm:"column:branchId" json:"branchId"`
	Branch            string    `gorm:"column:branch" json:"branch"`
	Plafond           float64   `gorm:"column:plafond" json:"plafond"`
	SubmittedLoanDate time.Time `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
	DisbursementDate  time.Time `gorm:"column:disbursementDate" json:"disbursementDate"`
}

type DisbursementStageInput struct {
	LastDisbursement     string    `gorm:"column:lastDisbursement" json:"lastDisbursement"`
	NextDisbursement     string    `gorm:"column:nextDisbursement" json:"nextDisbursement"`
	LastDisbursementDate time.Time `gorm:"column:lastDisbursementDate" json:"lastDisbursementDate"`
	NextDisbursementDate time.Time `gorm:"column:nextDisbursementDate" json:"nextDisbursementDate"`
	Remark               string    `gorm:"column:remark" json:"remark"`
}

func (disbursementStageInput *DisbursementStageInput) UpdateDateValue() {
	layout := "2006-01-02 15:04:05"
	disbursementStageInput.LastDisbursementDate, _ = time.Parse(layout, disbursementStageInput.LastDisbursement)
	disbursementStageInput.NextDisbursementDate, _ = time.Parse(layout, disbursementStageInput.NextDisbursement)
}

type DisbursementDetailByGroup struct {
	InvestorId       uint64    `gorm:"column:investorId" json:"investorId"`
	GroupID          uint64    `gorm:"column:groupId" json:"groupId"`
	GroupName        string    `gorm:"column:groupName" json:"groupName"`
	BranchName       string    `gorm:"column:branchName" json:"branchName"`
	Borrower         string    `gorm:"column:borrower" json:"borrower"`
	BorrowerNo       string    `gorm:"column:borrowerNo" json:"borrowerNo"`
	LoanID           uint64    `gorm:"column:loanId" json:"loanId"`
	Plafond          float64   `gorm:"column:plafond" json:"plafond"`
	DisbursementDate time.Time `gorm:"column:disbursementDate" json:"disbursementDate"`
	Stage            string    `gorm:"column:stage" json:"stage"`
	LoanStage        string    `gorm:"column:loanStage" json:"loanStage"`
}
