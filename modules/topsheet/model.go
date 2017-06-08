package topsheet

import "time"

type CurrentTopsheetSchema struct {
	LoanID              uint64    `gorm:"column:loanId" json:"loanId"`
	GroupID             uint64    `gorm:"column:groupId" json:"groupId"`
	GroupName           string    `gorm:"column:groupName" json:"groupName"`
	BranchID            uint64    `gorm:"column:branchId" json:"branchId"`
	ScheduleDay         string    `gorm:"column:scheduleDay" json:"scheduleDay"`
	ScheduleTime        string    `gorm:"column:scheduleTime" json:"scheduleTime"`
	Tenor               uint64    `gorm:"column:tenor" json:"tenor"`
	Rate                float64   `gorm:"column:rate" json:"rate"`
	Installment         float64   `gorm:"column:installment" json:"installment"`
	Plafond             float64   `gorm:"column:plafond" json:"plafond"`
	Subgroup            string    `gorm:"column:subgroup" json:"subgroup"`
	BorrowerNo          string    `gorm:"column:borrowerNo" json:"borrowerNo"`
	BorrowerName        string    `gorm:"column:borrowerName" json:"borrowerName"`
	Frequency           uint64    `gorm:"column:frequency" json:"frequency"`
	TotalReserve        float64   `gorm:"column:totalReserve" json:"totalReserve"`
	LatestReserve       float64   `gorm:"column:latestReserve" json:"latestReserve"`
	LatestInstallmentID uint64    `gorm:"column:latestInstallmentId" json:"latestInstallmentId"`
	TotalHadir          uint64    `gorm:"column:totalHadir" json:"totalHadir"`
	TotalAlfa           uint64    `gorm:"column:totalAlfa" json:"totalAlfa"`
	TotalCuti           uint64    `gorm:"column:totalCuti" json:"totalCuti"`
	TotalSakit          uint64    `gorm:"column:totalSakit" json:"totalSakit"`
	TotalIzin           uint64    `gorm:"column:totalIzin" json:"totalIzin"`
	SubmittedLoanDate   string    `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
	DisbursementDate    time.Time `gorm:"column:disbursementDate" json:"disbursementDate"`
}

type TopsheetFormSchema struct {
	LoanID          uint64    `gorm:"column:loanId" json:"loanId"`
	Type            string    `gorm:"column:type" json:"type"`
	Presence        string    `gorm:"column:presence" json:"presence"`
	PaidInstallment float64   `gorm:"column:paidInstallment" json:"paidInstallment"`
	Penalty         float64   `gorm:"column:penalty" json:"penalty"`
	Reserve         float64   `gorm:"column:reserve" json:"reserve"`
	Frequency       int32     `gorm:"column:frequency" json:"frequency"`
	Stage           string    `gorm:"column:stage" json:"stage"`
	CreatedAt       time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}
