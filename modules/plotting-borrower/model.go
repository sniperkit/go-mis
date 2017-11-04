package plottingBorrower

type EligbleInvestor struct {
	ID               uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Name             string `gorm:"column:name" json:"name"`
	InvestorNo       string `gorm:"column:investorNo" json:"investorNo"`
	BorrowerCriteria string `gorm:"column:borrowerCriteria" json:"borrowerCriteria"`
}

//TODO move to go-loan
type BorrowerCriteria struct{
	Area				[]Area		`json:"area"`
	Rate				Rate		`json:"rate"`
	Tenor				[]int		`json:"tenor"`
	Plafon				Plafon		`json:"plafond"`
	Sector				[]Sector		`json:"sector"`
	CreditScoreGrade	[]string	`json:"creditScoreGrade"`

}

//TODO move to go-loan
type Area struct{
	ID               int `json:"_id"`
	Name             string `json:"name"`
}

//TODO move to go-loan
type Rate struct{
	To               float64 `json:"to"`
	From             float64 `json:"from"`
	OptionType       uint64  `json:"optionType"`
}

//TODO move to go-loan
type Plafon struct{
	To               int `json:"to"`
	From             int `json:"from"`
	OptionType       uint64  `json:"optionType"`
}

//TODO move to go-loan
type Sector struct{
	ID               int `json:"_id"`
	Name             string `json:"name"`
}