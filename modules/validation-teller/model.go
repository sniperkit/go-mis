package validationTeller

type RawInstallmentDetail struct {
	Id            int64  `gorm:"column:id" json:"id"`
	BorrowerId            int64  `gorm:"column:borrowerId" json:"borrowerId"`
	Name               string  `gorm:"column:name" json:"name"`
	Repayment          float64 `gorm:"column:repayment" json:"repayment"`
	Tabungan           float64 `gorm:"column:tabungan" json:"tabungan"`
	Total              float64 `gorm:"column:total" json:"total"`
	Status			   string `gorm:"column:status" json:"status"`
	CashOnHand		   float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
	CashOnReserve		   float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
}

type RawInstallmentData struct {
	Fullname           string  `gorm:"column:fullname" json:"fullname"`
	GroupId            int64  `gorm:"column:groupId" json:"groupId"`
	Name               string  `gorm:"column:name" json:"name"`
	Repayment          float64 `gorm:"column:repayment" json:"repayment"`
	Tabungan           float64 `gorm:"column:tabungan" json:"tabungan"`
	Total              float64 `gorm:"column:total" json:"total"`
	TotalCair          float64 `gorm:"column:totalCair" json:"totalCair"`
	TotalGagalDropping float64 `gorm:"column:totalGagalDropping" json:"totalGagalDropping"`
	Status			   string `gorm:"column:status" json:"status"`
	CashOnHand		   float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
	CashOnReserve		   float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
}

type Majelis struct {
	GroupId            int64  `gorm:"column:groupId" json:"groupId"`
	Name               string  `gorm:"column:name" json:"name"`
	Repayment          float64 `gorm:"column:repayment" json:"repayment"`
	Tabungan           float64 `gorm:"column:tabungan" json:"tabungan"`
	Total              float64 `gorm:"column:total" json:"total"`
	TotalCair          float64 `gorm:"column:totalCair" json:"totalCair"`
	TotalGagalDropping float64 `gorm:"column:totalGagalDropping" json:"totalGagalDropping"`
	Status			   string `gorm:"column:status" json:"status"`
	CashOnHand		   float64 `gorm:"column:cashOnHand" json:"cashOnHand"`
	CashOnReserve		   float64 `gorm:"column:cashOnReserve" json:"cashOnReserve"`
}

type InstallmentData struct {
	Agent                string `gorm:"column:fullname" json:"fullname"`
	Majelis              []Majelis
	TotalActualRepayment float64 `gorm:"column:totalActualRepayment" json:"totalActualRepayment"`
}

// Coh - Cash on hand struct
type Coh struct {
	InstallmentId uint64
	cash          float64
}

// TellerValidation struct
type TellerValidation struct {
	ID         string `json:"id"`
	CashOnHand []Coh
}

// Log struct
type Log struct {
	GroupID   string      `json:"groupId"`
	ArchiveID string      `json:"archiveId"`
	Data      interface{} `json:"data"`
}