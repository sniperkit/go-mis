package cashout

import "time"

type Cashout struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CashoutID string     `gorm:"column:cashoutId" json:"cashoutId"`
	Amount    float64    `gorm:"column:amount" json:"amount"`
	Stage     string     `gorm:"column:stage" json:"stage"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// type CashoutInvestor struct {
// 	CashoutID    string    `gorm:"column:cashoutId" json:"cashoutId"`
// 	InvestorName string    `gorm:"column:investorName" json:"investorName"`
// 	Amount       float64   `gorm:"column:amount" json:"amount"`
// 	TotalBalance float64   `gorm:"column:totalBalance" json:"totalBalance"`
// 	CreatedAt    time.Time `gorm:"column:createdAt" json:"createdAt"`
// 	Stage        string    `gorm:"column:stage" json:"stage"`
// }

type CashoutInvestor struct {
	CifID           uint64     `gorm:"column:cifId" json:"cifId"`
	InvestorID      uint64     `gorm:"column:investorId" json:"investorId"`
	CashoutID       uint64     `gorm:"column:cashoutId" json:"cashoutId"`
	AccountID       uint64     `gorm:"column:accountId" json:"accountId"`
	InvestorName    string     `gorm:"column:investorName" json:"investorName"`
	CashoutNo       string     `gorm:"column:cashoutNo" json:"cashoutNo"`
	Amount          float64    `gorm:"column:amount" json:"amount"`
	TotalDebit      float64    `gorm:"column:totalDebit" json:"totalDebit"`
	TotalCredit     float64    `gorm:"column:totalCredit" json:"totalCredit"`
	TotalBalance    float64    `gorm:"column:totalBalance" json:"totalBalance"`
	Type            string     `gorm:"column:type" json:"type"`
	TransactionDate *time.Time `gorm:"column:transactionDate" json:"transactionDate"`
	Remark          string     `gorm:"column:remark" json:"remark"`
	Stage           string     `gorm:"column:stage" json:"stage"`
	RowsFullCount   int        `gorm:"column:full_count"`
	UpdatedAt       *time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}
