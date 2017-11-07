package model

import "time"

// InstallmentSchema - Installment schema
type InstallmentSchema struct {
	ID              uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Type            string     `gorm:"column:type" json:"type"`         // NORMAL, DOUBLE, PELUNASAN DINI
	Presence        string     `gorm:"column:presence" json:"presence"` // ATTEND, ABSENT, TR1, TR2, TR3
	PaidInstallment *float64    `gorm:"column:paidInstallment" json:"paidInstallment"`
	Penalty         *float64    `gorm:"column:penalty" json:"penalty"`
	Reserve         *float64    `gorm:"column:reserve" json:"reserve"`
	Frequency       *int32      `gorm:"column:frequency" json:"frequency"`
	Stage           string     `gorm:"column:stage" json:"stage"`
	Agentid         uint64     `gorm:"column:agentId" json:"agentId"`
	CashOnHand      float64 	`gorm:"column:cash_on_hand" json:"cashOnHand"`
	CashOnReserve   float64 	`gorm:"column:cash_on_reserve" json:"cashOnReserve"`
	CreatedAt       time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// InstallmentHistorySchema - Installment history schema
type InstallmentHistorySchema struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	StageFrom string     `gorm:"column:stageFrom" json:"stageFrom"`
	StageTo   string     `gorm:"column:stageTo" json:"stageTo"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// RInstallmentHistorySchema - Installment history relation schema
type RInstallmentHistorySchema struct {
	ID                   uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InstallmentID        uint64     `gorm:"column:installmentId" json:"installmentId"`
	InstallmentHistoryID uint64     `gorm:"column:installmentHistoryId" json:"installmentHistoryId"`
	CreatedAt            time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt            time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt            *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// RLoanInstallment - Loan installment relation schema
type RLoanInstallmentSchema struct {
	ID            uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanID        uint64     `gorm:"column:loanId" json:"loanId"`
	InstallmentID uint64     `gorm:"column:installmentId" json:"installmentId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// RLoanGroupSchema - Loan group relation schema
type RLoanGroupSchema struct {
	ID            uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanID        uint64     `gorm:"column:loanId" json:"loanId"`
	GroupId 	  uint64     `gorm:"column:groupId" json:"groupId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
