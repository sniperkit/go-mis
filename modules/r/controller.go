package r

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&RCifAccessToken{})
	services.BaseCrudInitWithDomain("r-cif-access-token", RCifAccessToken{}, []RCifAccessToken{})

	services.DBCPsql.AutoMigrate(&RCifBorrower{})
	services.BaseCrudInitWithDomain("r-cif-borrower", RCifBorrower{}, []RCifBorrower{})

	services.DBCPsql.AutoMigrate(&RCifInvestor{})
	services.BaseCrudInitWithDomain("r-cif-investor", RCifInvestor{}, []RCifInvestor{})

	services.DBCPsql.AutoMigrate(&RAccountTransactionCredit{})
	services.BaseCrudInitWithDomain("r-account-transaction-credit", RAccountTransactionCredit{}, []RAccountTransactionCredit{})

	services.DBCPsql.AutoMigrate(&RAccountTransactionDebit{})
	services.BaseCrudInitWithDomain("r-account-transaction-debit", RAccountTransactionDebit{}, []RAccountTransactionDebit{})

	services.DBCPsql.AutoMigrate(&RNotificationBorrower{})
	services.BaseCrudInitWithDomain("r-notification-borrower", RNotificationBorrower{}, []RNotificationBorrower{})

	services.DBCPsql.AutoMigrate(&RNotificationInvestor{})
	services.BaseCrudInitWithDomain("r-notification-investor", RNotificationInvestor{}, []RNotificationInvestor{})

	services.DBCPsql.AutoMigrate(&RAccountBorrower{})
	services.BaseCrudInitWithDomain("r-account-borrower", RAccountBorrower{}, []RAccountBorrower{})

	services.DBCPsql.AutoMigrate(&RAccountInvestor{})
	services.BaseCrudInitWithDomain("r-account-investor", RAccountInvestor{}, []RAccountInvestor{})

	services.DBCPsql.AutoMigrate(&RInvestorVirtualAccount{})
	services.BaseCrudInitWithDomain("r-investor-virtual-account", RInvestorVirtualAccount{}, []RInvestorVirtualAccount{})

	services.DBCPsql.AutoMigrate(&RVirtualAccountStatement{})
	services.BaseCrudInitWithDomain("r-virtual-account-statement", RVirtualAccountStatement{}, []RVirtualAccountStatement{})

	services.DBCPsql.AutoMigrate(&RLoanBorrower{})
	services.BaseCrudInitWithDomain("r-loan-borrower", RLoanBorrower{}, []RLoanBorrower{})

	services.DBCPsql.AutoMigrate(&RLoanAccountTransactionCredit{})
	services.BaseCrudInitWithDomain("r-loan-transaction-credit", RLoanAccountTransactionCredit{}, []RLoanAccountTransactionCredit{})

	services.DBCPsql.AutoMigrate(&RLoanHistory{})
	services.BaseCrudInitWithDomain("r-loan-history", RLoanHistory{}, []RLoanHistory{})

	services.DBCPsql.AutoMigrate(&RLoanArea{})
	services.BaseCrudInitWithDomain("r-loan-area", RLoanArea{}, []RLoanArea{})

	services.DBCPsql.AutoMigrate(&RLoanBranch{})
	services.BaseCrudInitWithDomain("r-loan-branch", RLoanBranch{}, []RLoanBranch{})

	services.DBCPsql.AutoMigrate(&RLoanMonitoring{})
	services.BaseCrudInitWithDomain("r-loan-monitoring", RLoanMonitoring{}, []RLoanMonitoring{})

	services.DBCPsql.AutoMigrate(&RLoanGroup{})
	services.BaseCrudInitWithDomain("r-loan-group", RLoanGroup{}, []RLoanGroup{})

	services.DBCPsql.AutoMigrate(&RLoanInstallment{})
	services.BaseCrudInitWithDomain("r-loan-installment", RLoanInstallment{}, []RLoanInstallment{})

	services.DBCPsql.AutoMigrate(&RLoanSector{})
	services.BaseCrudInitWithDomain("r-loan-sector", RLoanSector{}, []RLoanSector{})

	services.DBCPsql.AutoMigrate(&RLoanDisbursement{})
	services.BaseCrudInitWithDomain("r-loan-disbursement", RLoanDisbursement{}, []RLoanDisbursement{})

	services.DBCPsql.AutoMigrate(&RDisbursementHistory{})
	services.BaseCrudInitWithDomain("r-disbursement-history", RDisbursementHistory{}, []RDisbursementHistory{})

	services.DBCPsql.AutoMigrate(&RInstallmentHistory{})
	services.BaseCrudInitWithDomain("r-installment-history", RInstallmentHistory{}, []RInstallmentHistory{})

	services.DBCPsql.AutoMigrate(&RInstallmentAccountTransactionCredit{})
	services.BaseCrudInitWithDomain("r-installment-transaction-debit", RInstallmentAccountTransactionCredit{}, []RInstallmentAccountTransactionCredit{})

	services.DBCPsql.AutoMigrate(&RAreaBranch{})
	services.BaseCrudInitWithDomain("r-area-branch", RAreaBranch{}, []RAreaBranch{})

	services.DBCPsql.AutoMigrate(&RAreaUserMis{})
	services.BaseCrudInitWithDomain("r-area-user-mis", RAreaUserMis{}, []RAreaUserMis{})

	services.DBCPsql.AutoMigrate(&RBranchAgent{})
	services.BaseCrudInitWithDomain("r-branch-agent", RBranchAgent{}, []RBranchAgent{})

	services.DBCPsql.AutoMigrate(&RBranchUserMis{})
	services.BaseCrudInitWithDomain("r-branch-user-mis", RBranchUserMis{}, []RBranchUserMis{})

	services.DBCPsql.AutoMigrate(&RGroupAgent{})
	services.BaseCrudInitWithDomain("r-group-agent", RGroupAgent{}, []RGroupAgent{})

	services.DBCPsql.AutoMigrate(&RUserMisRole{})
	services.BaseCrudInitWithDomain("r-user-mis-role", RUserMisRole{}, []RUserMisRole{})

	services.DBCPsql.AutoMigrate(&RAdjustmentSubmittedBy{})
	services.BaseCrudInitWithDomain("r-adjustment-submitted-by", RAdjustmentSubmittedBy{}, []RAdjustmentSubmittedBy{})

	services.DBCPsql.AutoMigrate(&RAdjustmentApprovedBy{})
	services.BaseCrudInitWithDomain("r-adjustment-approved-by", RAdjustmentApprovedBy{}, []RAdjustmentApprovedBy{})

	services.DBCPsql.AutoMigrate(&RInvestorCashout{})
	services.BaseCrudInitWithDomain("r-investor-cashout", RInvestorCashout{}, []RInvestorCashout{})

	services.DBCPsql.AutoMigrate(&RCashoutHistory{})
	services.BaseCrudInitWithDomain("r-cashout-history", RCashoutHistory{}, []RCashoutHistory{})

	services.DBCPsql.AutoMigrate(&RUserMisAccessToken{})
	services.BaseCrudInitWithDomain("r-user-mis-access-token", RUserMisAccessToken{}, []RUserMisAccessToken{})

	services.DBCPsql.AutoMigrate(&RAgentAccessToken{})
	services.BaseCrudInitWithDomain("r-agent-access-token", RAgentAccessToken{}, []RAgentAccessToken{})

	services.DBCPsql.AutoMigrate(&RGroupBranch{})
	services.BaseCrudInitWithDomain("r-group-branch", RGroupBranch{}, []RGroupBranch{})

	services.DBCPsql.AutoMigrate(&RGroupBorrower{})
	services.BaseCrudInitWithDomain("r-group-borrower", RGroupBorrower{}, []RGroupBorrower{})

	services.DBCPsql.AutoMigrate(&RAccountTransactionDebitInstallment{})
	services.BaseCrudInitWithDomain("r-account-transaction-debit-installment", RAccountTransactionDebitInstallment{}, []RAccountTransactionDebitInstallment{})

	services.DBCPsql.AutoMigrate(&RAccountTransactionCreditLoan{})
	services.BaseCrudInitWithDomain("r-account-transaction-credit-loan", RAccountTransactionCreditLoan{}, []RAccountTransactionCreditLoan{})

	services.DBCPsql.AutoMigrate(&RInvestorProductPricingLoan{})
	services.BaseCrudInitWithDomain("r-investor-product-pricing-loan", RInvestorProductPricingLoan{}, []RInvestorProductPricingLoan{})

	services.DBCPsql.AutoMigrate(&RLoanOrder{})
	services.BaseCrudInitWithDomain("r-loan-order", RLoanOrder{}, []RLoanOrder{})

	services.DBCPsql.AutoMigrate(&RAccountTransactionCreditCashout{})
	services.BaseCrudInitWithDomain("r-account-transaction-credit-cashout", RAccountTransactionCreditCashout{}, []RAccountTransactionCreditCashout{})
}
