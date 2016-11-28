package r

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&R_CifInvestorBorrower{})
	services.DBCPsql.AutoMigrate(&R_BorrowerAccount{})
	services.DBCPsql.AutoMigrate(&R_InvestorBorrowerWallet{})
	services.DBCPsql.AutoMigrate(&R_WalletTransaction{})
	services.DBCPsql.AutoMigrate(&R_NotificationInvestorBorrower{})
	services.DBCPsql.AutoMigrate(&R_InvestorProductPricingLoan{})
	services.DBCPsql.AutoMigrate(&R_LoanDisbursement{})
	services.DBCPsql.AutoMigrate(&R_LoanSector{})
	services.DBCPsql.AutoMigrate(&R_LoanInstallment{})
	services.DBCPsql.AutoMigrate(&R_LoanGroup{})
	services.DBCPsql.AutoMigrate(&R_LoanCampaign{})
	services.DBCPsql.AutoMigrate(&R_LoanBranch{})
	services.DBCPsql.AutoMigrate(&R_LoanMonitoring{})
	services.DBCPsql.AutoMigrate(&R_GroupAgent{})
	services.DBCPsql.AutoMigrate(&R_AgentBranch{})
	services.DBCPsql.AutoMigrate(&R_AgentIncentive{})
	services.DBCPsql.AutoMigrate(&R_AgentBorrowerProspective{})

	services.BaseCrudInitWithDomain("cif-investor-borrower", R_CifInvestorBorrower{}, []R_CifInvestorBorrower{})
	services.BaseCrudInitWithDomain("borrower-account", R_BorrowerAccount{}, []R_BorrowerAccount{})
	services.BaseCrudInitWithDomain("investor-borrower-wallet", R_InvestorBorrowerWallet{}, []R_InvestorBorrowerWallet{})
	services.BaseCrudInitWithDomain("wallet-transaction", R_WalletTransaction{}, []R_WalletTransaction{})
	services.BaseCrudInitWithDomain("notification-investor-borrower", R_NotificationInvestorBorrower{}, []R_NotificationInvestorBorrower{})
	services.BaseCrudInitWithDomain("investor-product-pricing-loan", R_InvestorProductPricingLoan{}, []R_InvestorProductPricingLoan{})
	services.BaseCrudInitWithDomain("loan-disbursement", R_LoanDisbursement{}, []R_LoanDisbursement{})
	services.BaseCrudInitWithDomain("loan-sector", R_LoanSector{}, []R_LoanSector{})
	services.BaseCrudInitWithDomain("loan-installment", R_LoanInstallment{}, []R_LoanInstallment{})
	services.BaseCrudInitWithDomain("loan-group", R_LoanGroup{}, []R_LoanGroup{})
	services.BaseCrudInitWithDomain("loan-campaign", R_LoanCampaign{}, []R_LoanCampaign{})
	services.BaseCrudInitWithDomain("loan-branch", R_LoanBranch{}, []R_LoanBranch{})
	services.BaseCrudInitWithDomain("loan-monitoring", R_LoanMonitoring{}, []R_LoanMonitoring{})
	services.BaseCrudInitWithDomain("group-agent", R_GroupAgent{}, []R_GroupAgent{})
	services.BaseCrudInitWithDomain("agent-branch", R_AgentBranch{}, []R_AgentBranch{})
	services.BaseCrudInitWithDomain("agent-incentive", R_AgentIncentive{}, []R_AgentIncentive{})
	services.BaseCrudInitWithDomain("agent-borrower-propspective", R_AgentBorrowerProspective{}, []R_AgentBorrowerProspective{})
}