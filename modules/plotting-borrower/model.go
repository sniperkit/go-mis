package plottingBorrower

type EligbleInvestor struct {
	ID               uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name             string `gorm:"column:name" json:"name"`
	InvestorNo       string `gorm:"column:investorNo" json:"investorNo"`
	BorrowerCriteria string `gorm:"column:borrowerCriteria" json:"borrowerCriteria"`
}

type RecommendedLoan struct {
	LoanId           uint64  `gorm:"column:loanId" json:"loanId"`
	BorrowerName     string  `gorm:"column:borrowerName" json:"borrowerName"`
	Group            string  `gorm:"column:group" json:"group"`
	Branch           string  `gorm:"column:branch" json:"branch"`
	DisbursementDate string  `gorm:"column:disbursementDate" json:"disbursementDate"`
	Plafond          float64 `gorm:"column:plafond" json:"plafond"`
	Rate             float64 `gorm:"column:rate" json:"rate"`
	Tenor            uint64  `gorm:"column:tenor" json:"tenor"`
	CreditScoreGrade string  `gorm:"column:creditScoreGrade" json:"creditScoreGrade"`
	Purpose          string  `gorm:"column:purpose" json:"purpose"`
}

type GOLoanSuccessResponse struct {
	Status  int               `json:"status"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    []RecommendedLoan `json:"data"`
}