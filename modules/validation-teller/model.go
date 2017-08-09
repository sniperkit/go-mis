package validationTeller

type RawInstallmentDetail struct {
	Id            int64   `gorm:"column:id" json:"id"`
	BorrowerId    int64   `gorm:"column:borrowerId" json:"borrowerId"`
	Name          string  `gorm:"column:name" json:"name"`
	Repayment     float64 `gorm:"column:repayment" json:"repayment"`
	Tabungan      float64 `gorm:"column:tabungan" json:"tabungan"`
	Total         float64 `gorm:"column:total" json:"total"`
	Status        string  `gorm:"column:status" json:"status"`
	CashOnHand    float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
	CashOnReserve float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
}

type RawInstallmentData struct {
	Fullname           string  `gorm:"column:fullname" json:"fullname"`
	GroupId            int64   `gorm:"column:groupId" json:"groupId"`
	Name               string  `gorm:"column:name" json:"name"`
	Repayment          float64 `gorm:"column:repayment" json:"repayment"`
	Tabungan           float64 `gorm:"column:tabungan" json:"tabungan"`
	Total              float64 `gorm:"column:total" json:"total"`
	TotalCair          float64 `gorm:"column:totalCair" json:"totalCair"`
	TotalGagalDropping float64 `gorm:"column:totalGagalDropping" json:"totalGagalDropping"`
	Status             string  `gorm:"column:status" json:"status"`
	CashOnHand         float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
	CashOnReserve      float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
}

type Majelis struct {
	GroupId            int64   `gorm:"column:groupId" json:"groupId"`
	Name               string  `gorm:"column:name" json:"name"`
	Repayment          float64 `gorm:"column:repayment" json:"repayment"`
	Tabungan           float64 `gorm:"column:tabungan" json:"tabungan"`
	Total              float64 `gorm:"column:total" json:"total"`
	TotalCair          float64 `gorm:"column:totalCair" json:"totalCair"`
	TotalGagalDropping float64 `gorm:"column:totalGagalDropping" json:"totalGagalDropping"`
	Status             string  `gorm:"column:status" json:"status"`
	CashOnHand         float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
	CashOnReserve      float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
}

type ResponseGetData struct {
	InstallmentData      []InstallmentData `gorm:"column:installmentData" json:"installmentData"`
	TotalActualRepayment float64           `gorm:"column:totalActualRepayment" json:"totalActualRepayment"`
	TotalCashOnHand      float64           `gorm:"column:totalCashOnHand" json:"totalCashOnHand"`
	TotalTabungan        float64           `gorm:"column:totalTabungan" json:"totalTabungan"`
	TotalCashOnReserve   float64           `gorm:"column:totalCashOnReserve" json:"totalCashOnReserve"`
	TotalCair            float64           `gorm:"column:totalCair" json:"totalCair"`
	TotalGagalDroping    float64           `gorm:"column:totalGagalDroping" json:"totalGagalDroping"`
	BorrowerNotes        []interface{}     `json:"BorrowerNotes, omitempty"`
	MajelisNotes         []interface{}     `json:"majelisNotes, omitempty"`
}

type InstallmentData struct {
	Agent                string `gorm:"column:fullname" json:"fullname"`
	Majelis              []Majelis
	TotalActualRepayment float64 `gorm:"column:totalActualRepayment" json:"totalActualRepayment"`
	TotalCashOnHand      float64 `gorm:"column:totalCashOnHand" json:"totalCashOnHand"`
	TotalTabungan        float64 `gorm:"column:totalTabungan" json:"totalTabungan"`
	TotalCashOnReserve   float64 `gorm:"column:totalCashOnReserve" json:"totalCashOnReserve"`
	TotalCair            float64 `gorm:"column:totalCair" json:"totalCair"`
	TotalGagalDroping    float64 `gorm:"column:totalGagalDroping" json:"totalGagalDroping"`
}

// Coh - Cash on hand struct
type Coh struct {
	InstallmentId int64
	cash          float64
}

// TellerValidation struct
type TellerValidation struct {
	ID         string `json:"id"`
	CashOnHand []Coh
}

// Log struct
type Log struct {
	ID        string      `json:"id,omitempty"`
	GroupID   string      `json:"groupId,omitempty"`
	ArchiveID string      `json:"archiveId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// SubmitBody - struct
type SubmitBody struct {
	BranchID int64  `json:"branchId"`
	Date     string `json:"date"`
}

// DataLog - struct to store loging data archive / installment
type DataLog struct {
	Data interface{}
}
