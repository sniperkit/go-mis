package loanInsurance

type LoanInsuranceSchema struct {
  LoanID string `gorm:"column:loanId" json:"loanId"`
  BorrowerName string `gorm:"column:borrowerName" json:"borrowerName"`
  TotalPar uint32 `gorm:"column:totalPar" json:"totalPar"`
  TotalOtherType uint32 `gorm:"column:totalOtherType" json:"totalOtherType"`
  IsInsurance bool `gorm:"column:isInsurance" json:"isInsurance" sql:"default:false"`
	IsInsuranceRequested bool `gorm:"column:isInsuranceRequested" json:"isInsuranceRequested" sql:"default:false"`
	IsInsuranceRefund bool `gorm:"column:isInsuranceRefund" json:"isInsuranceRefund" sql:"default:false"`
}