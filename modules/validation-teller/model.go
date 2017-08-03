package validationTeller

type ValidationTellerData struct {
	InstallmentID      string  `gorm:"column:installmentId" json:"installmentId"`
	Fullname           string  `gorm:"column:fullname" json:"fullname"`
	Name               string  `gorm:"column:name" json:"name"`
	Repayment          float64 `gorm:"column:repayment" json:"repayment"`
	Tabungan           float64 `gorm:"column:tabungan" json:"tabungan"`
	Total              float64 `gorm:"column:total" json:"total"`
	TotalCair          float64 `gorm:"column:totalCair" json:"totalCair"`
	TotalGagalDropping float64 `gorm:"column:totalGagalDropping" json:"totalGagalDropping"`
}
