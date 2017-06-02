package r

import "time"

// Relation 'cif' to `access-token`
type RCifAccessToken struct {
	ID            uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CifId         uint64     `gorm:"column:cifId" json:"cifId"`
	AccessTokenId uint64     `gorm:"column:accessTokenId" json:"accessTokenId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'cif' to `borrower`
type RCifBorrower struct {
	ID         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CifId      uint64     `gorm:"column:cifId" json:"cifId"`
	BorrowerId uint64     `gorm:"column:borrowerId" json:"borrowerId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'cif' to `investor`
type RCifInvestor struct {
	ID         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CifId      uint64     `gorm:"column:cifId" json:"cifId"`
	InvestorId uint64     `gorm:"column:investorId" json:"investorId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'account' to `account-transaction-credit`
type RAccountTransactionCredit struct {
	ID                         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AccountId                  uint64     `gorm:"column:accountId" json:"accountId"`
	AccountTransactionCreditId uint64     `gorm:"column:accountTransactionCreditId" json:"accountTransactionCreditId"`
	CreatedAt                  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'account' to `account-transaction-credit`
type RAccountTransactionCreditLoan struct {
	ID                         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	OrderId                    uint64     `gorm:"column:orderId" json:"orderId"`
	LoanId                     uint64     `gorm:"column:loanId" json:"loanId"`
	AccountTransactionCreditId uint64     `gorm:"column:accountTransactionCreditId" json:"accountTransactionCreditId"`
	CreatedAt                  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'account' to `account-transaction-debit`
type RAccountTransactionDebit struct {
	ID                        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AccountId                 uint64     `gorm:"column:accountId" json:"accountId"`
	AccountTransactionDebitId uint64     `gorm:"column:accountTransactionDebitId" json:"accountTransactionDebitId"`
	CreatedAt                 time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                 time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                 *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `account-transaction-debit` to `installment`
type RAccountTransactionDebitInstallment struct {
	ID                        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AccountTransactionDebitId uint64     `gorm:"column:accountTransactionDebitId" json:"accountTransactionDebitId"`
	InstallmentId             uint64     `gorm:"column:installmentId" json:"installmentId"`
	CreatedAt                 time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                 time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                 *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `borrower`
type RLoanBorrower struct {
	ID         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId     uint64     `gorm:"column:loanId" json:"loanId"`
	BorrowerId uint64     `gorm:"column:borrowerId" json:"borrowerId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `account-transaction-credit`
type RLoanAccountTransactionCredit struct {
	ID                         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId                     uint64     `gorm:"column:loanId" json:"loanId"`
	AccountTransactionCreditId uint64     `gorm:"column:accountTransactionCreditId" json:"accountTransactionCreditId"`
	CreatedAt                  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `history`
type RLoanHistory struct {
	ID            uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId        uint64     `gorm:"column:loanId" json:"loanId"`
	LoanHistoryId uint64     `gorm:"column:loanHistoryId" json:"loanHistoryId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `area`
type RLoanArea struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint64     `gorm:"column:loanId" json:"loanId"`
	AreaId    uint64     `gorm:"column:areaId" json:"areaId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `branch`
type RLoanBranch struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint64     `gorm:"column:loanId" json:"loanId"`
	BranchId  uint64     `gorm:"column:branchId" json:"branchId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `loan-monitoring`
type RLoanMonitoring struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId           uint64     `gorm:"column:loanId" json:"loanId"`
	LoanMonitoringId uint64     `gorm:"column:loanMonitoringId" json:"loanMonitoringId"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `group`
type RLoanGroup struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint64     `gorm:"column:loanId" json:"loanId"`
	GroupId   uint64     `gorm:"column:groupId" json:"groupId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `installment`
type RLoanInstallment struct {
	ID            uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId        uint64     `gorm:"column:loanId" json:"loanId"`
	InstallmentId uint64     `gorm:"column:installmentId" json:"installmentId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `sector`
type RLoanSector struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint64     `gorm:"column:loanId" json:"loanId"`
	SectorId  uint64     `gorm:"column:sectorId" json:"sectorId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'loan' to `disbursement`
type RLoanDisbursement struct {
	ID             uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId         uint64     `gorm:"column:loanId" json:"loanId"`
	DisbursementId uint64     `gorm:"column:disbursementId" json:"disbursementId"`
	CreatedAt      time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `disbursement` to `disbursement-history`
type RDisbursementHistory struct {
	ID                    uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	DisbursementId        uint64     `gorm:"column:disbursementId" json:"disbursementId"`
	DisbursementHistoryId uint64     `gorm:"column:disbursementHistoryId" json:"disbursementHistoryId"`
	CreatedAt             time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt             time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt             *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'installment' to `installment-history`
type RInstallmentHistory struct {
	ID                   uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InstallmentId        uint64     `gorm:"column:installmentId" json:"installmentId"`
	InstallmentHistoryId uint64     `gorm:"column:installmentHistoryId" json:"installmentHistoryId"`
	CreatedAt            time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt            time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt            *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'installment' to `account-transaction-credit`
type RInstallmentAccountTransactionCredit struct {
	ID                        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InstallmentId             uint64     `gorm:"column:installmentId" json:"installmentId"`
	AccountTransactionDebitId uint64     `gorm:"column:accountTransactionDebitId" json:"accountTransactionDebitId"`
	CreatedAt                 time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                 time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                 *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'notification' to `borrower`
type RNotificationBorrower struct {
	ID             uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	NotificationId uint64     `gorm:"column:notificationId" json:"notificationId"`
	BorrowerId     uint64     `gorm:"column:borrowerId" json:"borrowerId"`
	CreatedAt      time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'notification' to `investor`
type RNotificationInvestor struct {
	ID             uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	NotificationId uint64     `gorm:"column:notificationId" json:"notificationId"`
	InvestorId     uint64     `gorm:"column:investorId" json:"investorId"`
	CreatedAt      time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'account' to `borrower`
type RAccountBorrower struct {
	ID         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AccountId  uint64     `gorm:"column:accountId" json:"accountId"`
	BorrowerId uint64     `gorm:"column:borrowerId" json:"borrowerId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'account' to `investor`
type RAccountInvestor struct {
	ID         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AccountId  uint64     `gorm:"column:accountId" json:"accountId"`
	InvestorId uint64     `gorm:"column:investorId" json:"investorId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'investor' to `virtual-account`
type RInvestorVirtualAccount struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InvestorId       uint64     `gorm:"column:investorId" json:"investorId"`
	VirtualAccountId uint64     `gorm:"column:vaId" json:"virtualAccountId"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'virtual-account' to `virtual-account-statement`
type RVirtualAccountStatement struct {
	ID                        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	VirtualAccountId          uint64     `gorm:"column:virtualAccountId" json:"virtualAccountId"`
	VirtualAccountStatementId uint64     `gorm:"column:virtualAccountStatementId" json:"virtualAccountStatementId"`
	CreatedAt                 time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                 time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                 *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'area' to `branch`
type RAreaBranch struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AreaId    uint64     `gorm:"column:areaId" json:"areaId"`
	BranchId  uint64     `gorm:"column:branchId" json:"branchId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'area' to `user-mis`
type RAreaUserMis struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AreaId    uint64     `gorm:"column:areaId" json:"areaId"`
	UserMisId uint64     `gorm:"column:userMisId" json:"userMisId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'branch' to `agent`
type RBranchAgent struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	BranchId  uint64     `gorm:"column:branchId" json:"branchId"`
	AgentId   uint64     `gorm:"column:agentId" json:"agentId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'group' to `agent`
type RGroupAgent struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	GroupId   uint64     `gorm:"column:groupId" json:"groupId"`
	AgentId   uint64     `gorm:"column:agentId" json:"agentId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation 'branch' to `user-mis`
type RBranchUserMis struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	BranchId  uint64     `gorm:"column:branchId" json:"branchId"`
	UserMisId uint64     `gorm:"column:userMisId" json:"userMisId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `user-mis` to 'role'
type RUserMisRole struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	UserMisId uint64     `gorm:"column:userMisId" json:"userMisId"`
	RoleId    uint64     `gorm:"column:roleId" json:"roleId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `adjustment` to 'user-mis' -> adjustment-submitted-by
type RAdjustmentSubmittedBy struct {
	ID           uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AdjustmentId uint64     `gorm:"column:adjustmentId" json:"adjustmentId"`
	UserMisId    uint64     `gorm:"column:userMisId" json:"userMisId"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `adjustment` to 'user-mis' -> adjustment-approved-by
type RAdjustmentApprovedBy struct {
	ID           uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AdjustmentId uint64     `gorm:"column:adjustmentId" json:"adjustmentId"`
	UserMisId    uint64     `gorm:"column:userMisId" json:"userMisId"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `investor` to 'cashout'
type RInvestorCashout struct {
	ID         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InvestorId uint64     `gorm:"column:investorId" json:"investorId"`
	CashoutId  uint64     `gorm:"column:cashoutId" json:"cashoutId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `cashout` to 'cashout-history'
type RCashoutHistory struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CashoutId        uint64     `gorm:"column:cashoutId" json:"cashoutId"`
	CashoutHistoryId uint64     `gorm:"column:cashoutHistoryId" json:"cashoutHistoryId"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `user-mis` to `access-token`
type RUserMisAccessToken struct {
	ID            uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	UserMisId     uint64     `gorm:"column:userMisId" json:"userMisId"`
	AccessTokenId uint64     `gorm:"column:accessTokenId" json:"accessTokenId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `agent` to `access-token`
type RAgentAccessToken struct {
	ID            uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AgentId       uint64     `gorm:"column:agentId" json:"agentId"`
	AccessTokenId uint64     `gorm:"column:accessTokenId" json:"accessTokenId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `group` to `branch`
type RGroupBranch struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	GroupId   uint64     `gorm:"column:groupId" json:"groupId"`
	BranchId  uint64     `gorm:"column:branchId" json:"branchId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `group` to `borrower`
type RGroupBorrower struct {
	ID         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	GroupId    uint64     `gorm:"column:groupId" json:"groupId"`
	BorrowerId uint64     `gorm:"column:borrowerId" json:"borrowerId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `investor` to `product-pricing` to `loan`
type RInvestorProductPricingLoan struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InvestorId       uint64     `gorm:"column:investorId" json:"investorId"`
	ProductPricingId uint64     `gorm:"column:productPricingId" json:"productPricingId"`
	LoanId           uint64     `gorm:"column:loanId" json:"loanId"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type RInvestorProductPricing struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InvestorId       uint64     `gorm:"column:investorId" json:"investorId"`
	ProductPricingId uint64     `gorm:"column:productPricingId" json:"productPricingId"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `loan-order`
type RLoanOrder struct {
	ID          uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanOrderId uint64     `gorm:"column:loanOrderId" json:"loanOrderId"`
	LoanId      uint64     `gorm:"column:loanId" json:"loanId"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `account_transaction_credit` to `cashout`
type RAccountTransactionCreditCashout struct {
	ID                         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AccountTransactionCreditID uint64     `gorm:"column:accountTransactionCreditId" json:"accountTransactionCreditId"`
	CashoutID                  uint64     `gorm:"column:cashoutId" json:"cashoutId"`
	CreatedAt                  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `adjustment` to `account_transaction_debit`
type RAdjustmentAccountTransactionDebit struct {
	ID                        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AdjustmentID              uint64     `gorm:"column:adjustmentId" json:"adjustmentId"`
	AccountTransactionDebitID uint64     `gorm:"column:accountTransactionDebitId" json:"accountTransactionDebitId"`
	CreatedAt                 time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt                 time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt                 *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `adjustment` to `installment`
type RInstallmentAdjustment struct {
	ID            uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AdjustmentID  uint64     `gorm:"column:adjustmentId" json:"adjustmentId"`
	InstallmentID uint64     `gorm:"column:installmentId" json:"installmentId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan_order` to `campaign`
type RLoanOrderCampaign struct {
	ID          uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanOrderID uint64     `gorm:"column:loanOrderId" json:"loanOrderId"`
	CampaignID  uint64     `gorm:"column:campaignId" json:"campaignId"`
	Quantity    uint64     `gorm:"column:quantity" json:"quantity"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
