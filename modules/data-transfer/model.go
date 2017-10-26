package dataTransfer

// DataTransfer - Data transfer table
type DataTransfer struct {
	ID                   uint64  `gorm:"column:id"`
	ValidationDate       string  `gorm:"column:validation_date" json:"validationDate"`
	TransferDate         string  `gorm:"column:transfer_date" json:"transferDate"`
	RepaymentID          string  `gorm:"column:repayment_id" json:"repaymentId"`
	RepaymentNominal     float64 `gorm:"column:repayment_nominal" json:"repaymentNominal"`
	RepaymentNote    string  `gorm:"column:repayment_note" json:"repaymentNote"`
	TabunganID           string  `gorm:"column:tabungan_id" json:"tabunganId"`
	TabunganNominal      float64 `gorm:"column:tabungan_nominal" json:"tabunganNominal"`
	TabunganNote    string  `gorm:"column:tabungan_note" json:"tabunganNote"`
	GagalDroppingID      string  `gorm:"column:gagal_dropping_id" json:"gagalDroppingId"`
	GagalDroppingNominal float64 `gorm:"column:gagal_dropping_nominal" json:"gagalDroppingNominal"`
	GagalDroppingNote    string  `gorm:"column:gagal_dropping_note" json:"gagalDroppingNote"`
	BranchID             uint64  `gorm:"column:branch_id" json:"branchId"`
}
