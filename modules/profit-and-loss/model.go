package profitAndLoss

import "time"

type ProfitAndLoss struct {
  ID uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
  LoanID uint64 `gorm:"column:loanId" json:"loanId"`
  Collector string `gorm:"column:collector" json:"collector"`
  Amount float64 `gorm:"column:amount" json:"amount"`
  CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}