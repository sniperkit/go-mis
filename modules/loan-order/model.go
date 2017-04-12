package loanOrder

import "time"

type LoanOrder struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	OrderNo   string     `gorm:"column:orderNo" json:"orderNo"`
	Remark    string     `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type LoanOrderInvestorPendingWaiting struct {
	ID           uint64  `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	OrderNo      string  `gorm:"column:orderNo" json:"orderNo"`
	Name         string  `gorm:"name:remark" json:"name"`
	TotalBalance float64 `gorm:"totalBalance:remark" json:"totalBalance"`
	TotalPlafond float64 `gorm:"totalPlafond:remark" json:"totalPlafond"`
}

type LoanOrderCompact struct {
	ID           uint64  `json:"_id"`
	Username     string  `json:"username"`
	Name         string  `json:"name"`
	OrderNo      string  `gorm:"column:orderNo" json:"orderNo"`
	TotalBalance float64 `gorm:"column:totalBalance" json:"totalBalance"`
	TotalPlafond float64 `gorm:"column:totalPlafond" json:"totalPlafond"`
	Remark       string  `json:"remark"`
}

type LoanOrderDetail struct {
	ID           uint64  `json:"_id"`
	Username     string  `json:"username"`
	Name         string  `json:"name"`
	OrderNo      string  `gorm:"column:orderNo" json:"orderNo"`
	LoanId       uint64  `gorm:"column:loanId" json:"loanId"`
	TotalBalance float64 `gorm:"column:totalBalance" json:"totalBalance"`
	Plafond 		 float64 `gorm:"column:plafond" json:"plafond"`
	Remark       string  `json:"remark"`			
}

type InvestorSearch struct {
	ID                uint64  `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Investor					uint64  `gorm:"column:loanId" json:"loanId"`
}
