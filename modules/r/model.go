package r

import "time"

// Relation `cif` to `investor` and `borrower`
type R_CifInvestorBorrower struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	CifId     uint64     `gorm:"column:cifId" json:"cifId"`
	DataId    uint64     `gorm:"column:dataId" json:"dataId"` // dataId: [ 'investorId', 'borrowerId' ]
	Type      string     `gorm:"column:type" json:"type"`     // type: [ 'investor', 'borrower' ]
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `borrower` to `account`
type R_BorrowerAccount struct {
	ID         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	BorrowerId uint64     `gorm:"column:borrowerId" json:"borrowerId"`
	AccountId  uint64     `gorm:"column:accountId" json:"accountId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `investor` and `borrower` to `wallet`
type R_InvestorBorrowerWallet struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	DataId    uint64     `gorm:"column:dataId" json:"dataId"` // dataId: [ 'investorId', 'borrowerId' ]
	WalletId  uint64     `gorm:"column:walletId" json:"walletId"`
	Type      string     `gorm:"column:type" json:"type"` // type: [ 'investor', 'borrower' ]
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `wallet` to `wallet-transaction`
type R_WalletTransaction struct {
	ID                  uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	WalletId            uint64     `gorm:"column:walletId" json:"walletId"`
	WalletTransactionId uint64     `gorm:"column:walletTransactionId" json:"walletTransactionId"`
	CreatedAt           time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt           *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `notification` to `investor` and `borrower`
type R_NotificationInvestorBorrower struct {
	ID             uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	NotificationId uint64     `gorm:"column:notificationId" json:"notificationId"`
	DataId         uint64     `gorm:"column:dataId" json:"dataId"` // dataId: [ 'investorId', 'borrowerId' ]
	Type           string     `gorm:"column:type" json:"type"`     // type: [ 'investor', 'borrower' ]
	CreatedAt      time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `investor` to `product-pricing` to `loan`
type R_InvestorProductPricingLoan struct {
	ID               uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InvestorId       uint64     `gorm:"column:investorId" json:"investorId"`
	ProductPricingId uint64     `gorm:"column:productPricingId" json:"productPricingId"`
	LoanId           uint64     `gorm:"column:loanId" json:"loanId"`
	CreatedAt        time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `disbursement`
type R_LoanDisbursement struct {
	ID             uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId         uint64     `gorm:"column:loanId" json:"loanId"`
	DisbursementId uint64     `gorm:"column:disbursementId" json:"disbursementId"`
	CreatedAt      time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `sector`
type R_LoanSector struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint64     `gorm:"column:loanId" json:"loanId"`
	SectorId  uint64     `gorm:"column:sectorId" json:"sectorId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `installment`
type R_LoanInstallment struct {
	ID            uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId        uint64     `gorm:"column:loanId" json:"loanId"`
	InstallmentId uint64     `gorm:"column:installmentId" json:"installmentId"`
	CreatedAt     time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `group`
type R_LoanGroup struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint64     `gorm:"column:loanId" json:"loanId"`
	GroupId   uint64     `gorm:"column:groupId" json:"groupId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `campaign`
type R_LoanCampaign struct {
	ID         uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId     uint64     `gorm:"column:loanId" json:"loanId"`
	CampaignId uint64     `gorm:"column:investorId" json:"investorId"`
	CreatedAt  time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `branch`
type R_LoanBranch struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint64     `gorm:"column:loanId" json:"loanId"`
	BranchId  uint64     `gorm:"column:branchId" json:"branchId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `loan` to `monitoring`
type R_LoanMonitoring struct {
	ID           uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId       uint64     `gorm:"column:loanId" json:"loanId"`
	MonitoringId uint64     `gorm:"column:monitoringId" json:"monitoringId"`
	CreatedAt    time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `installment` to `installment-presence`
type R_InstallmentPresence struct {
	ID                    uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	InstallmentId         uint64     `gorm:"column:installmentId" json:"installmentId"`
	InstallmentPresenceId uint64     `gorm:"column:installmentPresenceId" json:"installmentPresenceId"`
	CreatedAt             time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt             time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt             *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `group` to `agent`
type R_GroupAgent struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	GroupId   uint64     `gorm:"column:groupId" json:"groupId"`
	AgentId   uint64     `gorm:"column:agentId" json:"agentId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `agent` to `branch`
type R_AgentBranch struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AgentId   uint64     `gorm:"column:agentId" json:"agentId"`
	BranchId  uint64     `gorm:"column:branchId" json:"branchId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `agent` to `incentive`
type R_AgentIncentive struct {
	ID          uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AgentId     uint64     `gorm:"column:agentId" json:"agentId"`
	IncentiveId uint64     `gorm:"column:incentiveId" json:"incentiveId"`
	CreatedAt   time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `agent` to borrower-prospective
type R_AgentBorrowerProspective struct {
	ID                    uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	AgentId               uint64     `gorm:"column:agentInd" json:"agentId"`
	BorrowerProspectiveId uint64     `gorm:"column:borrowerProspectiveId" json:"borrowerProspectiveId"`
	CreatedAt             time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt             time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt             *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `user-mis` to `role`
type R_UserMisRole struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	UserMisId uint64     `gorm:"column:userMisId" json:"userMisId"`
	RoleId    uint64     `gorm:"column:roleId" json:"roleId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

// Relation `order` to `loan`
type R_OrderLoan struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	OrderId   uint64     `gorm:"column:orderId" json:"orderId"`
	LoanId    uint64     `gorm:"column:loanId" json:"loanId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}

type R_LoanAccount struct {
	ID        uint64     `gorm:"primary_key" gorm:"column:_id" json:"_id"`
	LoanId    uint64     `gorm:"column:loanId" json:"loanId"`
	AccountId uint64     `gorm:"column:accountId" json:"accountId"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
}
