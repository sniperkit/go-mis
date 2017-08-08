package installment

import (
	"time"
)

var query string

type Installment struct {
	ID              uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type            string     `gorm:"column:type" json:"type"`         // NORMAL, DOUBLE, PELUNASAN DINI
	Presence        string     `gorm:"column:presence" json:"presence"` // ATTEND, ABSENT, TR1, TR2, TR3
	PaidInstallment float64    `gorm:"column:paidInstallment" json:"paidInstallment"`
	Penalty         float64    `gorm:"column:penalty" json:"penalty"`
	Reserve         float64    `gorm:"column:reserve" json:"reserve"`
	Frequency       int32      `gorm:"column:frequency" json:"frequency"`
	Stage           string     `gorm:"column:stage" json:"stage"`
	CreatedAt       time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
	TransactionDate *time.Time `gorm:"column:transactionDate" json:"transactionDate"`
}

type InstallmentFetch struct {
	ID                   uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Branch               string    `gorm:"column:branch" json:"branch"`
	GroupID              uint64    `gorm:"column:groupId" json:"groupId"`
	Group                string    `gorm:"column:group" json:"group"`
	TotalPaidInstallment float64   `gorm:"column:totalPaidInstallment" json:"totalPaidInstallment"`
	TotalReserve         float64   `gorm:"column:totalReserve" json:"totalReserve"`
	CreatedAt            time.Time `gorm:"column:createdAt" json:"createdAt"`
}

type InstallmentDetail struct {
	GroupID         uint64  `gorm:"column:groupId" json:"groupId"`
	BranchID        uint64  `gorm:"column:branchId" json:"branchId"`
	GroupName       string  `gorm:"column:groupName" json:"groupName"`
	CifName         string  `gorm:"column:cifName" json:"cifName"`
	BorrowerNo      string  `gorm:"column:borrowerNo" json:"borrowerNo"`
	LoanId          string  `gorm:"column:loanId" json:"loanId"`
	InstallmentID   uint64  `gorm:"column:installmentId" json:"installmentId"`
	Type            string  `gorm:"column:type" json:"type"`
	Presence        string  `gorm:"column:presence" json:"presence"`
	PaidInstallment float64 `gorm:"column:paidInstallment" json:"paidInstallment"`
	Penalty         float64 `gorm:"column:penalty" json:"penalty"`
	Reserve         float64 `gorm:"column:reserve" json:"reserve"`
	Frequency       int32   `gorm:"column:frequency" json:"frequency"`
	Stage           string  `gorm:"column:stage" json:"stage"`
}
