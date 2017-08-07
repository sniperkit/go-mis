package validationTeller

type RawInstallmentData struct {
	Fullname           string  `gorm:"column:fullname" json:"fullname"`
	GroupId            string  `gorm:"column:groupId" json:"groupId"`
	Name               string  `gorm:"column:name" json:"name"`
	Repayment          float64 `gorm:"column:repayment" json:"repayment"`
	Tabungan           float64 `gorm:"column:tabungan" json:"tabungan"`
	Total              float64 `gorm:"column:total" json:"total"`
	TotalCair          float64 `gorm:"column:totalCair" json:"totalCair"`
	TotalGagalDropping float64 `gorm:"column:totalGagalDropping" json:"totalGagalDropping"`
	Status			   string `gorm:"column:status" json:"status"`
	CashOnHand		   string `gorm:"column:cashOnHand" json:"cashOnHand"`
	CashOnReserve		   string `gorm:"column:cashOnReserve" json:"cashOnReserve"`
}

type Majelis struct {
	GroupId            string  `gorm:"column:groupId" json:"groupId"`
	Name               string  `gorm:"column:name" json:"name"`
	Repayment          float64 `gorm:"column:repayment" json:"repayment"`
	Tabungan           float64 `gorm:"column:tabungan" json:"tabungan"`
	Total              float64 `gorm:"column:total" json:"total"`
	TotalCair          float64 `gorm:"column:totalCair" json:"totalCair"`
	TotalGagalDropping float64 `gorm:"column:totalGagalDropping" json:"totalGagalDropping"`
	Status			   string `gorm:"column:status" json:"status"`
	CashOnHand		   string `gorm:"column:cashOnHand" json:"cashOnHand"`
	CashOnReserve		   string `gorm:"column:cashOnReserve" json:"cashOnReserve"`
}

type InstallmentData struct {
	Agent                string `gorm:"column:fullname" json:"fullname"`
	Majelis              []Majelis
	TotalActualRepayment float64 `gorm:"column:totalActualRepayment" json:"totalActualRepayment"`
}
