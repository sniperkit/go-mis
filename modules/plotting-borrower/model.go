package plottingBorrower

type EligbleInvestor struct {
	ID               uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name             string `gorm:"column:name" json:"name"`
	InvestorNo       string `gorm:"column:investorNo" json:"investorNo"`
	BorrowerCriteria string `gorm:"column:borrowerCriteria" json:"borrowerCriteria"`
}
