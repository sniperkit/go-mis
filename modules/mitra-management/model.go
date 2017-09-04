package mitramanagement

type (
	Status struct {
		ID          uint64 `json:"id" gorm:"column:id"`
		Description string `json:"description" gorm:"column:description"`
		Type        string `json:"type" gorm:"column:type"`
	}

	Reason struct {
		ID          uint64 `json:"id" gorm:"column:id"`
		StatusID    uint64 `json:"statusId" gorm:"column:statusId"`
		Description string `json:"description" gorm:"description"`
	}

	// PortfolioAtRisk - Portfolio at risk information to the client
	PortfolioAtRisk struct {
		LoanID         uint64 `json:"loanId" gorm:"column:loanId"`
		BorrowerNumber uint64 `json:"borrowerNumber" gorm:"column:borrowerNumber"`
		BorrowerName   string `json:"borrowerName" gorm:"column:borrowerName"`
		MajelisName    string `json:"majelisName" gorm:"column:majelisName"`
		Reason         string `json:"reason" gorm:"column:reason"`
	}
)
