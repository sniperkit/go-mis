package dataTransfer

// DataTransfer - Data transfer table
type DataTransfer struct {
	ID                   uint64  `bson:"column:id"`
	ValidationDate       string  `bson:"column:validationDate" json:"validationDate"`
	TransferDate         string  `bson:"column:transferDate" json:"transferDate"`
	RepaymentID          string  `bson:"column:repaymentId" json:"repaymentId"`
	RepaymentNominal     float64 `bson:"column:repaymentNominal" json:"repaymentNominal"`
	TabunganID           string  `bson:"column:tabunganId" json:"tabunganId"`
	TabunganNominal      float64 `bson:"column:tabunganNominal" json:"tabunganNominal"`
	GagalDroppingID      string  `bson:"column:gagalDroppingId" json:"gagalDroppingId"`
	GagalDroppingNominal float64 `bson:"column:gagalDroppingNominal" json:"gagalDroppingNominal"`
	BranchID             uint64  `bson:"branchId" json:"branchId"`
}
