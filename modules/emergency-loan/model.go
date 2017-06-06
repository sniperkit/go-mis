package emergency_loan

type EmergencyLoanBorrower struct {
	ID           uint64     `gorm:"primary_key" gorm:"column:id" json:"id"`
	BorrowerId   uint64     `gorm:"column:borrowerId" json:"borrowerId"`
	BorrowerName string     `gorm:"column:borrowerName" json:"borrowerName"`
	BranchId     uint64     `gorm:"column:branchId" json:"branchId"`
	OldLoanId    uint64     `gorm:"column:oldLoanId" json:"oldLoanId"`
	NewLoanId    uint64     `gorm:"column:newLoanId" json:"newLoanId"`
	GroupId    uint64     `gorm:"column:groupId" json:"groupId"`
	Status       bool     `gorm:"column:status" json:"status"`
}
