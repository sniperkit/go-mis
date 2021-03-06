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

type InvestorDetail struct {
	InvestorId   uint64     `gorm:"column:investorId" json:"investorId"`
	Username   string     `gorm:"column:username" json:"username"`
	Name   string     `gorm:"column:name" json:"name"`
	PhoneNo   string     `gorm:"column:phoneNo" json:"phoneNo"`
}

type LoanOrderInvestorPendingWaiting struct {
	ID                     uint64  `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	OrderNo                string  `gorm:"column:orderNo" json:"orderNo"`
	Name                   string  `gorm:"name:remark" json:"name"`
	TotalBalance           float64 `gorm:"totalBalance:remark" json:"totalBalance"`
	TotalPlafond           float64 `gorm:"totalPlafond:remark" json:"totalPlafond"`
	UsingVoucher           bool    `gorm:"usingVoucher" json:"usingVoucher"`
	ParticipateCampaign    bool    `gorm:"column:participateCampaign" json:"participateCampaign"`
	QuantityOfCampaignItem uint64  `gorm:"column:quantityOfCampaignItem" json:"quantityOfCampaignItem"`
	VoucherAmount          float64 `gorm:"voucherAmount" json:"voucherAmount"`
}

type LoanOrderList struct {
	ID                     uint64    `gorm:"column:id" json:"_id"`
	Threshold               string    `gorm:"threshold" json:"threshold"`
	Username               string    `gorm:"username" json:"username"`
	Name                   string    `gorm:"name" json:"name"`
	OrderNo                string    `gorm:"column:orderNo" json:"orderNo"`
	TotalBalance           float64   `gorm:"column:totalBalance" json:"totalBalance"`
	TotalPlafond           float64   `gorm:"column:totalPlafond" json:"totalPlafond"`
	UsingVoucher           bool      `gorm:"column:usingVoucher" json:"usingVoucher"`
	VoucherAmount          float64   `gorm:"column:voucherAmount" json:"voucherAmount"`
	InsuranceAmount          float64 `gorm:"column:insuranceAmount" json:"insuranceAmount"`
	ParticipateCampaign    bool      `gorm:"column:participateCampaign" json:"participateCampaign"`
	QuantityOfCampaignItem uint64    `gorm:"column:quantityOfCampaignItem" json:"quantityOfCampaignItem"`
	CampaignAmount         float64   `gorm:"column:campaignAmount" json:"campaignAmount"`
	UsingRefreal    bool      `gorm:"column:usingRefreal" json:"usingRefreal"`
	CreatedAt              time.Time `gorm:"column:createdAt" json:"createdAt"`
}

type LoanOrderDetail struct {
	ID                     uint64  `gorm:"column:id" json:"_id"`
	Username               string  `gorm:"username" json:"username"`
	Threshold               string    `gorm:"threshold" json:"threshold"`
	Name                   string  `gorm:"name" json:"name"`
	OrderNo                string  `gorm:"column:orderNo" json:"orderNo"`
	LoanId                 uint64  `gorm:"column:loanId" json:"loanId"`
	Purpose                string  `gorm:"column:purpose" json:"purpose"`
	TotalBalance           float64 `gorm:"column:totalBalance" json:"totalBalance"`
	UsingVoucher           bool    `gorm:"column:usingVoucher" json:"usingVoucher"`
	VoucherAmount          float64 `gorm:"column:voucherAmount" json:"voucherAmount"`
	InsuranceAmount          float64 `gorm:"column:insuranceAmount" json:"insuranceAmount"`
	ParticipateCampaign    bool    `gorm:"column:participateCampaign" json:"participateCampaign"`
	QuantityOfCampaignItem uint64  `gorm:"column:quantityOfCampaignItem" json:"quantityOfCampaignItem"`
	CampaignAmount         float64 `gorm:"column:campaignAmount" json:"campaignAmount"`
	Plafond                float64 `gorm:"column:plafond" json:"plafond"`
	UsingRefreal    bool      `gorm:"column:usingRefreal" json:"usingRefreal"`
	Remark                 string  `gorm:"column:remark" json:"remark"`
}

type InvestorSearch struct {
	ID       uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Investor uint64 `gorm:"column:loanId" json:"loanId"`
}
