package validationTeller

import "bitbucket.org/go-mis/modules/installment"

type (
	RawInstallmentDetail struct {
		Id                int64   `gorm:"column:id" json:"id"`
		BorrowerId        int64   `gorm:"column:borrowerId" json:"borrowerId"`
		Name              string  `gorm:"column:name" json:"name"`
		Repayment         float64 `gorm:"column:repayment" json:"repayment"`
		Tabungan          float64 `gorm:"column:tabungan" json:"tabungan"`
		Total             float64 `gorm:"column:total" json:"total"`
		Status            string  `gorm:"column:status" json:"status"`
		CashOnHand        float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
		CashOnReserve     float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
		CashOnHandNote    string  `gorm:"column:cashOnHandNote" json:"cashOnHandNote"`
		CashOnReserveNote string  `gorm:"column:cashOnReserveNote" json:"cashOnReserveNote"`
	}

	Notes struct {
		GroupId       int64   `gorm:"column:groupId" json:"groupId"`
		Name          string  `gorm:"column:name" json:"name"`
		CashOnHand    float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
		CashOnReserve float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
		Disbursement float64 `gorm:"column:disbursement" json:"disbursement"`
		Note          string  `gorm:"column:note" json:"note"`
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

	MajelisId struct {
		GroupId int64  `gorm:"column:groupId" json:"groupId"`
		Name    string `gorm:"column:name" json:"name"`
	}

	ResponseGetData struct {
		InstallmentData      []InstallmentData `gorm:"column:installmentData" json:"installmentData"`
		TotalActualRepayment float64           `gorm:"column:totalActualRepayment" json:"totalActualRepayment"`
		TotalCashOnHand      float64           `gorm:"column:totalCashOnHand" json:"totalCashOnHand"`
		TotalTabungan        float64           `gorm:"column:totalTabungan" json:"totalTabungan"`
		TotalCashOnReserve   float64           `gorm:"column:totalCashOnReserve" json:"totalCashOnReserve"`
		TotalCair            float64           `gorm:"column:totalCair" json:"totalCair"`
		TotalGagalDroping    float64           `gorm:"column:totalGagalDroping" json:"totalGagalDroping"`
		BorrowerNotes        interface{}       `json:"borrowerNotes, omitempty"`
		MajelisNotes         interface{}       `json:"majelisNotes, omitempty"`
		ListMajelis          []MajelisId       `json:"listMajelis, omitempty"`
		IsEnableSubmit       bool              `json:"isEnableSubmit"`
		DataTransfer         DataTransfer      `json:"dataTransfer,omitempty"`
	}

	InstallmentData struct {
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

	// Coh - Cash on hand struct
	Coh struct {
		InstallmentId int64
		cash          float64
	}

	// TellerValidation struct
	TellerValidation struct {
		ID         string `json:"id"`
		CashOnHand []Coh
	}

	// Log struct
	Log struct {
		ID        string      `json:"id,omitempty"`
		GroupID   string      `json:"groupId,omitempty"`
		ArchiveID string      `json:"archiveId,omitempty"`
		Data      interface{} `json:"data,omitempty"`
	}

	// SubmitBody - struct
	SubmitBody struct {
		BranchID int64  `json:"branchId"`
		Date     string `json:"date"`
	}

	// DataLog - struct to store loging data archive / installment
	DataLog struct {
		Data interface{}
	}

	DataTransfer struct {
		ID                   uint64  `json:"id,omitempty"`
		ValidationDate       string  `json:"validationDate"`
		TransferDate         string  `json:"transferDate"`
		RepaymentID          string  `json:"repaymentId"`
		RepaymentNominal     float64 `json:"repaymentNominal"`
		TabunganID           string  `json:"tabunganId"`
		TabunganNominal      float64 `json:"tabunganNominal"`
		GagalDroppingID      string  `json:"gagalDroppingId"`
		GagalDroppingNominal float64 `json:"gagalDroppingNominal"`
		GagalDroppingNote string `json:"gagalDroppingNote"`
		BranchID             uint64  `json:"branchId"`
	}

	TotalCabang struct {
		TotalCabRepaymentAct       float64
		TotalCabRepaymentCoh       float64
		TotalCabTabunganAct        float64
		TotalCabTabunganCoh        float64
		TotalCabActualAgent        float64
		TotalCabCohAgent           float64
		TotalCabPencairanAgent     float64
		TotalCabGagalDroppingAgent float64
	}

	TotalRepayment struct {
		TotalRepaymentAct       float64
		TotalRepaymentProj      float64
		TotalRepaymentCoh       float64
		TotalTabunganAct        float64
		TotalTabunganProj       float64
		TotalTabunganCoh        float64
		TotalActualAgent        float64
		TotalProjectionAgent    float64
		TotalCohAgent           float64
		TotalPencairanAgent     float64
		TotalPencairanProjAgent float64
		TotalGagalDroppingAgent float64
	}
)

func (totalRepayment *TotalRepayment) AddTotal(rawInstallmentData installment.RawInstallmentData) {
	totalRepayment.TotalRepaymentAct += rawInstallmentData.Repayment
	totalRepayment.TotalRepaymentProj += rawInstallmentData.ProjectionRepayment
	totalRepayment.TotalRepaymentCoh += rawInstallmentData.CashOnHand
	totalRepayment.TotalTabunganAct += rawInstallmentData.Tabungan
	totalRepayment.TotalTabunganProj += rawInstallmentData.ProjectionTabungan
	totalRepayment.TotalTabunganCoh += rawInstallmentData.CashOnReserve
	totalRepayment.TotalActualAgent += rawInstallmentData.Total
	totalRepayment.TotalProjectionAgent += rawInstallmentData.ProjectionRepayment + rawInstallmentData.ProjectionTabungan
	totalRepayment.TotalCohAgent += rawInstallmentData.CashOnHand + rawInstallmentData.CashOnReserve
	totalRepayment.TotalPencairanAgent += rawInstallmentData.TotalCair
	totalRepayment.TotalPencairanProjAgent += rawInstallmentData.TotalCairProj
	totalRepayment.TotalGagalDroppingAgent += rawInstallmentData.TotalGagalDropping
}

func (totalCabang *TotalCabang) AddTotal(totalRepayment *TotalRepayment) {
	totalCabang.TotalCabRepaymentAct += totalRepayment.TotalRepaymentAct
	totalCabang.TotalCabRepaymentCoh += totalRepayment.TotalRepaymentCoh
	totalCabang.TotalCabTabunganAct += totalRepayment.TotalTabunganAct
	totalCabang.TotalCabTabunganCoh += totalRepayment.TotalTabunganCoh
	totalCabang.TotalCabActualAgent += totalRepayment.TotalActualAgent
	totalCabang.TotalCabCohAgent += totalRepayment.TotalCohAgent
	totalCabang.TotalCabPencairanAgent += totalRepayment.TotalPencairanAgent
	totalCabang.TotalCabGagalDroppingAgent += totalRepayment.TotalGagalDroppingAgent
}

func (installmentData *InstallmentData) AddTotal(totalRepayment *TotalRepayment) {
	installmentData.TotalActualRepayment = totalRepayment.TotalRepaymentAct
	installmentData.TotalProjectionRepayment = totalRepayment.TotalRepaymentProj
	installmentData.TotalCohRepayment = totalRepayment.TotalRepaymentCoh
	installmentData.TotalActualTabungan = totalRepayment.TotalTabunganAct
	installmentData.TotalProjectionTabungan = totalRepayment.TotalTabunganProj
	installmentData.TotalCohTabungan = totalRepayment.TotalTabunganCoh
	installmentData.TotalActualAgent = totalRepayment.TotalActualAgent
	installmentData.TotalProjectionAgent = totalRepayment.TotalProjectionAgent
	installmentData.TotalCohAgent = totalRepayment.TotalCohAgent
	installmentData.TotalPencairanAgent = totalRepayment.TotalPencairanAgent
	installmentData.TotalPencairanProjAgent = totalRepayment.TotalPencairanProjAgent
	installmentData.TotalGagalDroppingAgent = totalRepayment.TotalGagalDroppingAgent
}

func (majelis Majelis) InitializedByRawInstallmentData(rawInstallmentData installment.RawInstallmentData) Majelis {
	m := Majelis{
		GroupId:             rawInstallmentData.GroupId,
		Name:                rawInstallmentData.Name,
		Repayment:           rawInstallmentData.Repayment,
		Tabungan:            rawInstallmentData.Tabungan,
		TotalActual:         rawInstallmentData.Total,
		TotalProyeksi:       rawInstallmentData.ProjectionRepayment + rawInstallmentData.ProjectionTabungan,
		TotalCoh:            rawInstallmentData.CashOnHand + rawInstallmentData.CashOnReserve,
		TotalCair:           rawInstallmentData.TotalCair,
		TotalCairProj:       rawInstallmentData.TotalCairProj,
		TotalGagalDropping:  rawInstallmentData.TotalGagalDropping,
		Status:              rawInstallmentData.Status,
		CashOnHand:          rawInstallmentData.CashOnHand,
		CashOnReserve:       rawInstallmentData.CashOnReserve,
		ProjectionRepayment: rawInstallmentData.ProjectionRepayment,
		ProjectionTabungan:  rawInstallmentData.ProjectionTabungan,
	}
	return m
}

func AssignTotalResponseData(responseData *ResponseGetData, totalCabang *TotalCabang) {
	responseData.TotalActualRepayment = totalCabang.TotalCabRepaymentAct
	responseData.TotalCashOnHand = totalCabang.TotalCabRepaymentCoh
	responseData.TotalCashOnReserve = totalCabang.TotalCabTabunganCoh
	responseData.TotalGagalDroping = totalCabang.TotalCabGagalDroppingAgent
	responseData.TotalTabungan = totalCabang.TotalCabTabunganAct
	responseData.TotalCair = totalCabang.TotalCabPencairanAgent
}
