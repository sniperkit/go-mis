package plottingBorrower

type EligbleInvestor struct {
	ID               uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name             string `gorm:"column:name" json:"name"`
	InvestorNo       string `gorm:"column:investorNo" json:"investorNo"`
	BorrowerCriteria string `gorm:"column:borrowerCriteria" json:"borrowerCriteria"`
}

//TODO move to go-loan
type BorrowerCriteria struct {
	Area             []Area   `json:"area"`
	Rate             Rate     `json:"rate"`
	Tenor            []int    `json:"tenor"`
	Plafon           Plafon   `json:"plafond"`
	Sector           []Sector `json:"sector"`
	CreditScoreGrade []string `json:"creditScoreGrade"`
}

//TODO move to go-loan
type Area struct {
	ID   int    `json:"_id"`
	Name string `json:"name"`
}

//TODO move to go-loan
type Rate struct {
	To         float64 `json:"to"`
	From       float64 `json:"from"`
	OptionType uint64  `json:"optionType"`
}

//TODO move to go-loan
type Plafon struct {
	To         int    `json:"to"`
	From       int    `json:"from"`
	OptionType uint64 `json:"optionType"`
}

//TODO move to go-loan
type Sector struct {
	ID   int    `json:"_id"`
	Name string `json:"name"`
}

type RecomendedLoan struct {
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
	Status  int              `json:"status"`
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    []RecomendedLoan `json:"data"`
}
