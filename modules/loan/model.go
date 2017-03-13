package loan

import "time"

type Loan struct {
	ID                   uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanPeriod           int64      `gorm:"column:loanPeriod" json:"loanPeriod"`
	AgreementType        string     `gorm:"column:agreementType" json:"agreementType"`
	Subgroup             string     `gorm:"column:subgroup" json:"subgrop"`
	Purpose              string     `gorm:"column:purpose" json:"purpose"`
	URLPic1              string     `gorm:"column:urlPic1" json:"urlPic1"`
	URLPic2              string     `gorm:"column:urlPic2" json:"urlPic2"`
	SubmittedLoanDate    string     `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
	SubmittedPlafond     float64    `gorm:"column:submittedPlafond" json:"submittedPlafond"`
	SubmittedTenor       int64      `gorm:"column:submittedTenor" json:"submittedTenor"`
	SubmittedInstallment float64    `gorm:"column:submittedInstallment" json:"submittedInstallment"`
	CreditScoreGrade     string     `gorm:"column:creditScoreGrade" json:"creditScoreGrade"`
	CreditScoreValue     float64    `gorm:"column:creditScoreValue" json:"creditScoreValue"`
	Tenor                uint64     `gorm:"column:tenor" json:"tenor"`
	Rate                 float64    `gorm:"column:rate" json:"rate"`
	Installment          float64    `gorm:"column:installment" json:"installment"`
	Plafond              float64    `gorm:"column:plafond" json:"plafond"`
	GroupReserve         float64    `gorm:"column:groupReserve" json:"groupReserve"`
	Stage                string     `gorm:"column:stage" json:"stage"`
	IsLWK                bool       `gorm:"column:isLWK" json:"isLWK" sql:"default:false"`
	IsUPK                bool       `gorm:"column:isUPK" json:"IsUPK" sql:"default:false"`
	CreatedAt            time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt            time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt            *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type LoanDatatable struct {
	ID                uint64  `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Borrower          string  `gorm:"column:borrower" json:"borrower"`
	Group             string  `gorm:"column:group" json:"group"`
	SubmittedLoanDate string  `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
	DisbursementDate  string  `gorm:"column:disbursementDate" json:"disbursementDate"`
	Plafond           float64 `gorm:"column:plafond" json:"plafond"`
	Tenor             uint64  `gorm:"column:tenor" json:"tenor"`
	Rate              float64 `gorm:"column:rate" json:"rate"`
	Stage             string  `gorm:"column:stage" json:"stage"`
}

type LoanFetch struct {
	ID                uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanPeriod        int64     `gorm:"column:loanPeriod" json:"loanPeriod"`
	AgreementType     string    `gorm:"column:agreementType" json:"agreementType"`
	Subgroup          string    `gorm:"column:subgroup" json:"subgrop"`
	SubmittedLoanDate time.Time `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
	CreditScoreGrade  string    `gorm:"column:creditScoreGrade" json:"creditScoreGrade"`
	CreditScoreValue  float64   `gorm:"column:creditScoreValue" json:"creditScoreValue"`
	Tenor             uint64    `gorm:"column:tenor" json:"tenor"`
	Rate              float64   `gorm:"column:rate" json:"rate"`
	Installment       float64   `gorm:"column:installment" json:"installment"`
	Plafond           float64   `gorm:"column:plafond" json:"plafond"`
	Stage             string    `gorm:"column:stage" json:"stage"`
	CreatedAt         time.Time `gorm:"column:createdAt" json:"createdAt"`
	Sector            string    `gorm:"column:sector" json:"sector"`
	Borrower          string    `gorm:"column:borrower" json:"borrower"`
	Group             string    `gorm:"column:group" json:"group"`
	Branch            string    `gorm:"column:branch" json:"branch"`
	DisbursementDate  string    `gorm:"column:disbursementDate" json:"disbursementDate"`
}

type LoanBorrowerProfile struct {
	CifNumber uint64 `gorm:"column:cifNumber" json:"cifNumber"`
	Name      string `gorm:"name" json:"borrower"`
	Area      string `gorm:"area" json:"area"`
	Branch    string `gorm:"branch" json:"branch"`
	Group     string `gorm:"group" json:"group"`
}

type Akad struct {
	ID            uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AgreementType string    `gorm:"column:agreementType" json:"agreementType"`
	Purpose       string    `gorm:"column:purpose" json:"purpose"`
	Plafond       float64   `gorm:"column:plafond" json:"plafond"`
	Tenor         uint64    `gorm:"column:tenor" json:"tenor"`
	Installment   float64   `gorm:"column:installment" json:"installment"`
	Investor      string    `gorm:"column:investor" json:"investor"`
	Borrower      string    `gorm:"column:borrower" json:"borrower"`
	Group         string    `gorm:"column:group" json:"group"`
	UpdatedAt     time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

type LoanDropping struct {
	ID       uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	Borrower string `gorm:"column:borrower" json:"borrower"`
	Group    string `gorm:"column:group" json:"group"`
	Stage    string `gorm:"column:stage" json:"stage"`
	Investor string `gorm:"column:investor" json:"investor"`
}

type RefundBase struct {
	LoanID     uint64  `gorm:"column:loan_id"`
	InvestorID uint64  `gorm:"column:investor_id"`
	AccountID  uint64  `gorm:"column:account_id"`
	Plafond    float64 `gorm:"column:plafond"`
}

type AccountSum struct {
	Sum float64 `gorm:"column:sum"`
}
