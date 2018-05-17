package dataTransfer


// DataTransfer - Data transfer table
type DataTransfer struct {
	ID                   uint64  `gorm:"column:id"`
	BranchID             uint64  `gorm:"column:branchId" json:"branchId"`
	TransactionType      string  `gorm:"column:transactionType" json:"transactionType"`
	TransferDate         string  `gorm:"column:transferDate" json:"transferDate"`
	ValidationDate       string  `gorm:"column:validationDate" json:"validationDate"`
    SettlementID         uint64 `gorm:"column:settlementId" json:"settlementId"`
    BankStatementID      uint64 `gorm:"column:bankStatementId" json:"bankStatementId"`
    Note                 string `gorm:"column:note" json:"note"`
    BankStatementMatched bool `gorm:"column:bankStatementMatched" json:"bankStatementMatched"`
    Amount float64 `gorm:"column:amount" json:"amount`
    ReferenceCode string `gorm:"column:referenceCode" json:"referenceCode"`
    TransferNoteURL string `gorm:"column:transferNoteUrl" json:"transferNoteUrl"`
}

type DataTransfers struct {
    Items []DataTransfer `json:"data"`
}

type ValidationDate struct {
    CreatedAt string `gorm:"column:createdAt"`
}
