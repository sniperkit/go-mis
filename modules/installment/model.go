package installment

import (
	"time"
)

type (
	Installment struct {
		ID                uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
		Type              string     `gorm:"column:type" json:"type"`         // NORMAL, DOUBLE, PELUNASAN DINI
		Presence          string     `gorm:"column:presence" json:"presence"` // ATTEND, ABSENT, TR1, TR2, TR3
		PaidInstallment   float64    `gorm:"column:paidInstallment" json:"paidInstallment"`
		Penalty           float64    `gorm:"column:penalty" json:"penalty"`
		Reserve           float64    `gorm:"column:reserve" json:"reserve"`
		Frequency         int32      `gorm:"column:frequency" json:"frequency"`
		Stage             string     `gorm:"column:stage" json:"stage"`
		CreatedAt         time.Time  `gorm:"column:createdAt" json:"createdAt"`
		UpdatedAt         time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
		DeletedAt         *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
		TransactionDate   *time.Time `gorm:"column:transactionDate" json:"transactionDate"`
		CashOnHandNote    string     `gorm:"column:cash_on_hand_note" json:"cashOnHandNote"`
		CashOnReserveNote string     `gorm:"column:cash_on_reserve_note" json:"cashOnReserveNote"`
		StatusID          uint64     `gorm:"statusId" json:"statusId"`
		ReasonID          uint64     `gorm:"reasonId" json:"reasonId"`
	}

	InstallmentFetch struct {
		ID                   uint64    `gorm:"primary_key" gorm:"column:_id" json:"_id"`
		Branch               string    `gorm:"column:branch" json:"branch"`
		GroupID              uint64    `gorm:"column:groupId" json:"groupId"`
		Group                string    `gorm:"column:group" json:"group"`
		TotalPaidInstallment float64   `gorm:"column:totalPaidInstallment" json:"totalPaidInstallment"`
		TotalReserve         float64   `gorm:"column:totalReserve" json:"totalReserve"`
		CreatedAt            time.Time `gorm:"column:createdAt" json:"createdAt"`
	}

	InstallmentDetail struct {
		InstallmentID       uint64  `gorm:"column:installmentId" json:"installmentId"`
		BorrowerID          string  `gorm:"column:borrowerId" json:"borrowerId"`
		BorrowerName        string  `gorm:"column:borrowerName" json:"borrowerName"`
		Type                string  `gorm:"column:type" json:"type"`
		Presence            string  `gorm:"column:presence" json:"presence"`
		Frequency           int32   `gorm:"column:frequency" json:"frequency"`
		Repayment           float64 `gorm:"column:repayment" json:"repayment"`
		Tabungan            float64 `gorm:"column:tabungan" json:"tabungan"`
		Total               float64 `gorm:"column:total" json:"total"`
		ProjectionRepayment float64 `gorm:"column:projectionRepayment" json:"projectionRepayment"`
		ProjectionTabungan  float64 `gorm:"column:projectionTabungan" json:"projectionTabungan"`
		PaidInstallment     float64 `gorm:"column:paidInstallment" json:"paidInstallment"`
		Reserve             float64 `gorm:"column:reserve" json:"reserve"`
		Stage               string  `gorm:"column:stage" json:"stage"`
		TotalCair           float64 `gorm:"column:totalCair" json:"totalCair"`
		CashOnHand          float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
		CashOnReserve       float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
		CashOnHandNote      string  `gorm:"column:cashOnHandNote" json:"cashOnHandNote"`
		CasOnReserveNote    string  `gorm:"column:cashOnReserveNote" json:"cashOnReserveNote"`
	}

	RawInstallmentData struct {
		Fullname            string  `gorm:"column:fullname" json:"fullname"`
		GroupId             int64   `gorm:"column:groupId" json:"groupId"`
		Name                string  `gorm:"column:name" json:"name"`
		Repayment           float64 `gorm:"column:repayment" json:"repayment"`
		Tabungan            float64 `gorm:"column:tabungan" json:"tabungan"`
		Total               float64 `gorm:"column:total" json:"total"`
		TotalCair           float64 `gorm:"column:totalCair" json:"totalCair"`
		TotalCairProj       float64 `gorm:"column:totalCairProj" json:"totalCairProj"`
		TotalGagalDropping  float64 `gorm:"column:totalGagalDropping" json:"totalGagalDropping"`
		Status              string  `gorm:"column:status" json:"status"`
		CashOnHand          float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
		CashOnReserve       float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
		ProjectionRepayment float64 `gorm:"column:projectionRepayment" json:"projectionRepayment"`
		ProjectionTabungan  float64 `gorm:"column:projectionTabungan" json:"projectionTabungan"`
	}

	Majelis struct {
		GroupId             int64   `gorm:"column:groupId" json:"groupId"`
		Name                string  `gorm:"column:name" json:"name"`
		Repayment           float64 `gorm:"column:repayment" json:"repayment"`
		Tabungan            float64 `gorm:"column:tabungan" json:"tabungan"`
		TotalActual         float64 `gorm:"column:totalActual" json:"totalActual"`
		TotalProyeksi       float64 `gorm:"column:totalProyeksi" json:"totalProyeksi"`
		TotalCoh            float64 `gorm:"column:totalCoh" json:"totalCoh"`
		TotalCair           float64 `gorm:"column:totalCair" json:"totalCair"`
		TotalCairProj       float64 `gorm:"column:totalCairProj" json:"totalCairProj"`
		TotalGagalDropping  float64 `gorm:"column:totalGagalDropping" json:"totalGagalDropping"`
		Status              string  `gorm:"column:status" json:"status"`
		CashOnHand          float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
		CashOnReserve       float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
		ProjectionRepayment float64 `gorm:"column:projectionRepayment" json:"projectionRepayment"`
		ProjectionTabungan  float64 `gorm:"column:projectionTabungan" json:"projectionTabungan"`
	}

	PendingInstallmentData struct {
		Agent                    string `gorm:"column:fullname" json:"fullname"`
		Majelis                  []Majelis
		TotalActualRepayment     float64 `gorm:"column:totalActualRepayment" json:"totalActualRepayment"`
		TotalActualTabungan      float64 `gorm:"column:totalActualTabungan" json:"totalActualTabungan"`
		TotalActualAgent         float64 `gorm:"column:totalActualAgent" json:"totalActualAgent"`
		TotalProjectionRepayment float64 `gorm:"column:totalProjectionRepayment" json:"totalProjectionRepayment"`
		TotalProjectionTabungan  float64 `gorm:"column:totalProjectionTabungan" json:"totalProjectionTabungan"`
		TotalProjectionAgent     float64 `gorm:"column:totalProjectionAgent" json:"totalProjectionAgent"`
		TotalCohRepayment        float64 `gorm:"column:totalCohRepayment" json:"totalCohRepayment"`
		TotalCohTabungan         float64 `gorm:"column:totalCohTabungan" json:"totalCohTabungan"`
		TotalCohAgent            float64 `gorm:"column:totalCohAgent" json:"totalCohAgent"`
		TotalPencairanAgent      float64 `gorm:"column:totalPencairanAgent" json:"totalPencairanAgent"`
		TotalPencairanProjAgent  float64 `gorm:"column:totalPencairanProjAgent" json:"totalPencairanProjAgent"`
		TotalGagalDroppingAgent  float64 `gorm:"column:totalGagalDroppingAgent" json:"totalGagalDroppingAgent"`
	}

	PendingInstallment struct {
		PendingInstallmentData []PendingInstallmentData `json:"pendingInstallmentData, omitempty"`
		BorrowerNotes          interface{}              `json:"borrowerNotes, omitempty"`
		MajelisNotes           interface{}              `json:"majelisNotes, omitempty"`
		ListMajelis            []MajelisId              `json:"listMajelis,omitempty"`
	}

	MajelisId struct {
		GroupId int64  `gorm:"column:groupId" json:"groupId"`
		Name    string `gorm:"column:name" json:"name"`
	}
	LoanInvestorAccountID struct {
		LoanID     uint64  `gorm:"column:loanId" json:"loanId"`
		InvestorID uint64  `gorm:"column:investorId" json:"investorId"`
		AccountID  uint64  `gorm:"column:accountId" json:"accountId"`
		PPLROI     float64 `gorm:"column:pplROI" json:"pplROI"`
	}

	AccountTransactionDebitAndCredit struct {
		TotalDebit  float64 `gorm:"column:totalDebit" json:"totalDebit"`
		TotalCredit float64 `gorm:"column:totalCredit" json:"totalCredit"`
	}

	LoanSchema struct {
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
)
