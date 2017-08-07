package installment

import (
	"errors"
	"log"
	"strings"
	"time"

	"bitbucket.org/go-mis/services"
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

// FindByBranchAndDate - Filter Installment by branch and date
func (i Installment) FindByBranchAndDate(branchID, transactionDate string) ([]Installment, error) {
	if len(strings.Trim(branchID, " ")) == 0 {
		return nil, errors.New("Branch ID can not be empty")
	}
	if len(strings.Trim(transactionDate, " ")) == 0 {
		return nil, errors.New("Transaction date can not be empty")
	}
	installemnts := make([]Installment, 0)
	db := services.DBCPsql.Begin()
	defer db.Close()
	query = `select installment.id,
				installment.type,
				installment.presence,
				installment."paidInstallment", 
				installment.penalty,
				installment.reserve,
				installment.frequency,
				installment.stage,
				installment."transactionDate"
			FROM installment,
				r_loan_installment,
				loan,
				branch,
				r_loan_branch
			WHERE installment.id = r_loan_installment."installmentId" AND 
			loan.id = r_loan_installment."loanId" AND
			loan.id = r_loan_branch."loanId" AND 
			branch.id = r_loan_branch."branchId" AND 
			installment."deletedAt" is null AND 
			UPPER(installment.stage) = 'TELLER'
			branch.id = ? AND
			installment."transactionDate" = ?`
	db.Raw(query, branchID, transactionDate).Scan(&installemnts)
	if db.Error != nil {
		log.Println("#ERROR: ", db.Error)
		return nil, errors.New("Unable to retrieve installments")
	}
	return installemnts, nil
}
