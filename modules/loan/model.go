package loan

import "time"

type Loan struct {
	ID                   uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanType           string      `gorm:"column:loanType" json:"loanType" sql:"default:NORMAL"`
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
	IsInsurance					 bool				`gorm:"column:isInsurance" json:"isInsurance" sql:"default:false"`
	IsInsuranceRequested bool 			`gorm:"column:isInsuranceRequested" json:"isInsuranceRequested" sql:"default:false"`
	IsInsuranceRefund		 bool 			`gorm:"column:isInsuranceRefund" json:"isInsuranceRefund" sql:"default:false"`
	InsuranceType				 string			`gorm:"column:insuranceType" json:"insuranceType"` // JAMKRINDO, ALLIANZ
	EndType							 string     `gorm:"column:endType" json:"endType"` // EARLY, PENDING, PELDIN, INSURANCE
	CreatedAt            time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt            time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt            *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type LoanDatatable struct {
	ID                uint64  `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanID            uint64  `gorm:"column:loanId" json:"loanId"`
	BorrowerNo        string  `gorm:"column:borrowerNo" json:"borrowerNo"`
	InvestorId       	uint64  `gorm:"column:investorId" json:"investorId"`
	Borrower          string  `gorm:"column:borrower" json:"borrower"`
	Group             string  `gorm:"column:group" json:"group"`
	SubmittedLoanDate string  `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
	DisbursementDate  string  `gorm:"column:disbursementDate" json:"disbursementDate"`
	Plafond           float64 `gorm:"column:plafond" json:"plafond"`
	Tenor             uint64  `gorm:"column:tenor" json:"tenor"`
	Rate              float64 `gorm:"column:rate" json:"rate"`
	Stage             string  `gorm:"column:stage" json:"stage"`
	AkadAvailable     bool    `gorm:"column:akadAvailable" json:"akadAvailable"`
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
	BorrowerID	uint64    `gorm:"column:borrowerId" json:"borrowerId"`
	CifNumber uint64 `gorm:"column:cifNumber" json:"cifNumber"`
	Name      string `gorm:"name" json:"borrower"`
	Area      string `gorm:"area" json:"area"`
	Branch    string `gorm:"branch" json:"branch"`
	Group     string `gorm:"group" json:"group"`
}

type Akad struct {
	ID                 uint64  `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AgreementType      string  `gorm:"column:agreementType" json:"agreementType"`
	Purpose            string  `gorm:"column:purpose" json:"purpose"`
	Plafond            float64 `gorm:"column:plafond" json:"plafond"`
	Tenor              uint64  `gorm:"column:tenor" json:"tenor"`
	Installment        float64 `gorm:"column:installment" json:"installment"`
	Rate               float64 `gorm:"column:rate" json:"rate"`
	SubmittedLoanDate  string  `gorm:"column:submittedLoanDate" json:"submittedLoanDate"`
	InvestorID         uint64  `gorm:"column:investorId" json:"investorId"`
	Investor           string  `gorm:"column:investor" json:"investor"`
	Borrower           string  `gorm:"column:borrower" json:"borrower"`
	Group              string  `gorm:"column:group" json:"group"`
	ReturnOfInvestment float64 `gorm:"column:returnOfInvestment" json:"returnOfInvestment"`
	AdminitrationFee   float64 `gorm:"column:administrationFee" json:"administrationFee"`
	ServiceFee         float64 `gorm:"column:serviceFee" json:"serviceFee"`
	DisbursementDate   string  `gorm:"column:disbursementDate" json:"disbursementDate"`
	AgreementNo		   string  `gorm:"column:agreementNo" json:"agreementNo"`
}

type LoanDropping struct {
	ID         uint64 `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	BorrowerNo string `gorm:"column:borrowerNo" json:"borrowerNo"`
	Borrower   string `gorm:"column:borrower" json:"borrower"`
	Group      string `gorm:"column:group" json:"group"`
	Stage      string `gorm:"column:stage" json:"stage"`
	InvestorID uint64 `gorm:"column:investorId" json:"investorId"`
	Investor   string `gorm:"column:investor" json:"investor"`
}

type LoanStageHistory struct {
	ID        uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	StageFrom string    `gorm:"column:stageFrom" json:"stageFrom"`
	StageTo   string    `gorm:"column:stageTo" json:"stageTo"`
	Remark    string    `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
}

type RefundBase struct {
	LoanID     uint64  `gorm:"column:loan_id"`
	InvestorID uint64  `gorm:"column:investor_id"`
	AccountID  uint64  `gorm:"column:account_id"`
	Plafond    float64 `gorm:"column:plafond"`
	IsInsurance bool `gorm:"column:isInsurance"`
}

type AccountSum struct {
	Sum float64 `gorm:"column:sum"`
}
