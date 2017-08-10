package installment

import (
	"time"
)

type (
	Installment struct {
		ID              uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
		Type            string     `gorm:"column:type" json:"type"`         // NORMAL, DOUBLE, PELUNASAN DINI
		Presence        string     `gorm:"column:presence" json:"presence"` // ATTEND, ABSENT, TR1, TR2, TR3
		PaidInstallment float64    `gorm:"column:paidInstallment" json:"paidInstallment"`
		Penalty         float64    `gorm:"column:penalty" json:"penalty"`
		Reserve         float64    `gorm:"column:reserve" json:"reserve"`
		Frequency       int32      `gorm:"column:frequency" json:"frequency"`
		Stage           string     `gorm:"column:stage" json:"stage"`
		CreatedAt       time.Time  `gorm:"column:createdAt" json:"createdAt"`
		UpdatedAt       time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
		DeletedAt       *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
		TransactionDate *time.Time `gorm:"column:transactionDate" json:"transactionDate"`
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
		GroupID         uint64  `gorm:"column:groupId" json:"groupId"`
		BranchID        uint64  `gorm:"column:branchId" json:"branchId"`
		GroupName       string  `gorm:"column:groupName" json:"groupName"`
		CifName         string  `gorm:"column:cifName" json:"cifName"`
		BorrowerNo      string  `gorm:"column:borrowerNo" json:"borrowerNo"`
		LoanId          string  `gorm:"column:loanId" json:"loanId"`
		InstallmentID   uint64  `gorm:"column:installmentId" json:"installmentId"`
		Type            string  `gorm:"column:type" json:"type"`
		Presence        string  `gorm:"column:presence" json:"presence"`
		PaidInstallment float64 `gorm:"column:paidInstallment" json:"paidInstallment"`
		Penalty         float64 `gorm:"column:penalty" json:"penalty"`
		Reserve         float64 `gorm:"column:reserve" json:"reserve"`
		Frequency       int32   `gorm:"column:frequency" json:"frequency"`
		Stage           string  `gorm:"column:stage" json:"stage"`
	}

	PendingRawInstallmentData struct {
		Fullname            string  `gorm:"column:fullname" json:"fullname"`
		GroupId             int64   `gorm:"column:groupId" json:"groupId"`
		Name                string  `gorm:"column:name" json:"name"`
		Repayment           float64 `gorm:"column:repayment" json:"repayment"`
		Tabungan            float64 `gorm:"column:tabungan" json:"tabungan"`
		Total               float64 `gorm:"column:total" json:"total"`
		TotalCair           float64 `gorm:"column:totalCair" json:"totalCair"`
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
		TotalGagalDroppingAgent  float64 `gorm:"column:totalGagalDroppingAgent" json:"totalGagalDroppingAgent"`
	}

	PendingInstallment struct {
		PendingInstallmentData []PendingInstallmentData `json:"pendingInstallmentData, omitempty"`
		BorrowerNotes          interface{}              `json:"borrowerNotes, omitempty"`
		MajelisNotes           interface{}              `json:"majelisNotes, omitempty"`
	}
)
