package mitramanagement

type (
	Status struct {
		ID          uint64 `json:"id" gorm:"column:id"`
		Description string `json:"description" gorm:"column:description"`
		Type        string `json:"type" gorm:"column:type"`
	}

	Reason struct {
		ID          uint64 `json:"id" gorm:"column:id"`
		StatusID    uint64 `json:"statusId,omitempty" gorm:"column:statusId"`
		Description string `json:"description" gorm:"column:description"`
	}

	MMBorrower struct {
		LoanID         uint64 `json:"loanId" gorm:"column:loanId"`
		InstallmentID  uint64 `json:"installmentID" gorm:"column:installmentId"`
		BorrowerNumber string `json:"borrowerNumber" gorm:"column:borrowerNumber"`
		BorrowerName   string `json:"borrowerName" gorm:"column:borrowerName"`
		GroupName      string `json:"groupName" gorm:"column:groupName"`
		Reason         string `json:"reason" gorm:"column:reason"`
	}

	MMDOBorrower struct {
		LoanID         uint64  `json:"loanId" gorm:"column:loanId"`
		InstallmentID  uint64  `json:"installmentID" gorm:"column:installmentId"`
		BorrowerNumber uint64  `json:"borrowerNumber" gorm:"column:borrowerNumber"`
		BorrowerName   string  `json:"borrowerName" gorm:"column:borrowerName"`
		GroupName      string  `json:"groupName" gorm:"column:groupName"`
		Reason         string  `json:"reason" gorm:"column:reason"`
		Plafond        float64 `json:"plafond" gorm:"column:plafond"`
		Tenor          uint64  `json:"tenor" gorm:"column:tenor"`
		DODate         string  `json:"doDate" gorm:"column:doDate"`
		Agent          string  `json:"agent" gorm:"column:agent"`
		Type           string  `json:"type" gorm:"column:type"`
	}

	MMPARBorrower struct {
		LoanID         uint64  `json:"loanId" gorm:"column:loanId"`
		InstallmentID  uint64  `json:"installmentID" gorm:"column:installmentId"`
		BorrowerNumber uint64  `json:"borrowerNumber" gorm:"column:borrowerNumber"`
		BorrowerName   string  `json:"borrowerName" gorm:"column:borrowerName"`
		GroupName      string  `json:"groupName" gorm:"column:groupName"`
		Reason         string  `json:"reason" gorm:"column:reason"`
		Nominal        float64 `json:"nominal" gorm:"column:nominal"`
		PARDate        string  `json:"parDate" gorm:"column:parDate"`
		Agent          string  `json:"agent" gorm:"column:agent"`
	}

	MMTRBorrower struct {
		LoanID          uint64  `json:"loanId" gorm:"column:loanId"`
		InstallmentID   uint64  `json:"installmentID" gorm:"column:installmentId"`
		BorrowerNumber  uint64  `json:"borrowerNumber" gorm:"column:borrowerNumber"`
		BorrowerName    string  `json:"borrowerName" gorm:"column:borrowerName"`
		GroupName       string  `json:"groupName" gorm:"column:groupName"`
		Reason          string  `json:"reason" gorm:"column:reason"`
		Nominal         float64 `json:"nominal" gorm:"column:nominal"`
		TRDate          string  `json:"parDate" gorm:"column:trDate"`
		Agent           string  `json:"agent" gorm:"column:agent"`
		InstallmentType string  `json:"installmentType" gorm:"column:type"`
	}
)
