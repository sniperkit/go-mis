package dataTransfer

type DataTransfer struct {
	ID                   uint64  `bson:"column:id"`
	ValidationDate       string  `bson:"column:validationDate" json:"validationDate"`
	TransferDate         string  `bson:"column:transferDate" json:"transferDate"`
	RepaymentID          uint64  `bson:"column:repaymentId" json:"repaymentId"`
	RepaymentNominal     float64 `bson:"column:repaymentNominal" json:"repaymentNominal"`
	TabunganID           uint64  `bson:"column:tabunganId" json:"tabunganId"`
	TabunganNominal      float64 `bson:"column:tabunganNominal" json:"tabunganNominal"`
	GagalDroppingID      uint64  `bson:"column:gagalDroppingId" json:"gagalDroppingId"`
	GagalDroppingNominal float64 `bson:"column:gagalDroppingNominal" json:"gagalDroppingNominal"`
}
