package r

import "time"

// Relation `cif` to `investor` and `borrower`
type R_CifInvestorBorrower struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CifId     uint       `gorm:"column:cifId" json:"cifId"`
	DataId    uint       `gorm:"column:dataId" json:"dataId"` // dataId: [ 'investorId', 'borrowerId' ]
	Type      string     `gorm:"column:type" json:"type"`     // type: [ 'investor', 'borrower' ]
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `borrower` to `account`
type R_BorrowerAccount struct {
	ID         uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	BorrowerId uint       `gorm:"column:borrowerId" json:"borrowerId"`
	AccountId  uint       `gorm:"column:accountId" json:"accountId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `investor` and `borrower` to `wallet`
type R_InvestorBorrowerWallet struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	DataId    uint       `gorm:"column:dataId" json:"dataId"` // dataId: [ 'investorId', 'borrowerId' ]
	WalletId  uint       `gorm:"column:walletId" json:"walletId"`
	Type      string     `gorm:"column:type" json:"type"` // type: [ 'investor', 'borrower' ]
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `wallet` to `wallet-transaction`
type R_WalletTransaction struct {
	ID                  uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	WalletId            uint       `gorm:"column:walletId" json:"walletId"`
	WalletTransactionId uint       `gorm:"column:walletTransactionId" json:"walletTransactionId"`
	CreatedAt           time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt           *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `notification` to `investor` and `borrower`
type R_NotificationInvestorBorrower struct {
	ID             uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	NotificationId uint       `gorm:"column:notificationId" json:"notificationId"`
	DataId         uint       `gorm:"column:dataId" json:"dataId"` // dataId: [ 'investorId', 'borrowerId' ]
	Type           string     `gorm:"column:type" json:"type"`     // type: [ 'investor', 'borrower' ]
	CreatedAt      time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `investor` to `product-pricing` to `loan`
type R_InvestorProductPricingLoan struct {
	ID               uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InvestorId       uint       `gorm:"column:investorId" json:"investorId"`
	ProductPricingId uint       `gorm:"column:productPricingId" json:"productPricingId"`
	LoanId           uint       `gorm:"column:loanId" json:"loanId"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `disbursement`
type R_LoanDisbursement struct {
	ID             uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId         uint       `gorm:"column:loanId" json:"loanId"`
	DisbursementId uint       `gorm:"column:disbursementId" json:"disbursementId"`
	CreatedAt      time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `sector`
type R_LoanSector struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint       `gorm:"column:loanId" json:"loanId"`
	SectorId  uint       `gorm:"column:sectorId" json:"sectorId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `installment`
type R_LoanInstallment struct {
	ID            uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId        uint       `gorm:"column:loanId" json:"loanId"`
	InstallmentId uint       `gorm:"column:installmentId" json:"installmentId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `group`
type R_LoanGroup struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint       `gorm:"column:loanId" json:"loanId"`
	GroupId   uint       `gorm:"column:groupId" json:"groupId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `campaign`
type R_LoanCampaign struct {
	ID         uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId     uint       `gorm:"column:loanId" json:"loanId"`
	CampaignId uint       `gorm:"column:investorId" json:"investorId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `branch`
type R_LoanBranch struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint       `gorm:"column:loanId" json:"loanId"`
	BranchId  uint       `gorm:"column:branchId" json:"branchId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `monitoring`
type R_LoanMonitoring struct {
	ID           uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId       uint       `gorm:"column:loanId" json:"loanId"`
	MonitoringId uint       `gorm:"column:monitoringId" json:"monitoringId"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `installment` to `installment-presence`
type R_InstallmentPresence struct {
	ID                    uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InstallmentId         uint       `gorm:"column:installmentId" json:"installmentId"`
	InstallmentPresenceId uint       `gorm:"column:installmentPresenceId" json:"installmentPresenceId"`
	CreatedAt             time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt             time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt             *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `group` to `agent`
type R_GroupAgent struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	GroupId   uint       `gorm:"column:groupId" json:"groupId"`
	AgentId   uint       `gorm:"column:agentId" json:"agentId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `agent` to `branch`
type R_AgentBranch struct {
	ID        uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AgentId   uint       `gorm:"column:agentId" json:"agentId"`
	BranchId  uint       `gorm:"column:branchId" json:"branchId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `agent` to `incentive`
type R_AgentIncentive struct {
	ID          uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AgentId     uint       `gorm:"column:agentId" json:"agentId"`
	IncentiveId uint       `gorm:"column:incentiveId" json:"incentiveId"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `agent` to borrower-prospective
type R_AgentBorrowerProspective struct {
	ID                    uint       `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AgentId               uint       `gorm:"column:agentInd" json:"agentId"`
	BorrowerProspectiveId uint       `gorm:"column:borrowerProspectiveId" json:"borrowerProspectiveId"`
	CreatedAt             time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt             time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt             *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
