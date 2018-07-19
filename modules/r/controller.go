package r

import (
	"fmt"

	"bitbucket.org/go-mis/config"
	loanRaw "bitbucket.org/go-mis/modules/loan-raw"
	"bitbucket.org/go-mis/services"
)

func Init() {
	if config.AutoMigrate {
		fmt.Println("AutoMigrate All Model R Controller")
		services.DBCPsql.AutoMigrate(&RCifAccessToken{})
		services.DBCPsql.AutoMigrate(&RCifBorrower{})
		services.DBCPsql.AutoMigrate(&RCifInvestor{})
		services.DBCPsql.AutoMigrate(&RAccountTransactionCredit{})
		services.DBCPsql.AutoMigrate(&RAccountTransactionDebit{})
		services.DBCPsql.AutoMigrate(&RNotificationBorrower{})
		services.DBCPsql.AutoMigrate(&RNotificationInvestor{})
		services.DBCPsql.AutoMigrate(&RAccountBorrower{})
		services.DBCPsql.AutoMigrate(&RAccountInvestor{})
		services.DBCPsql.AutoMigrate(&RInvestorVirtualAccount{})
		services.DBCPsql.AutoMigrate(&RVirtualAccountStatement{})
		services.DBCPsql.AutoMigrate(&RLoanBorrower{})
		services.DBCPsql.AutoMigrate(&RLoanAccountTransactionCredit{})
		services.DBCPsql.AutoMigrate(&RLoanHistory{})
		services.DBCPsql.AutoMigrate(&RLoanArea{})
		services.DBCPsql.AutoMigrate(&RLoanBranch{})
		services.DBCPsql.AutoMigrate(&RLoanMonitoring{})
		services.DBCPsql.AutoMigrate(&RLoanGroup{})
		services.DBCPsql.AutoMigrate(&RLoanInstallment{})
		services.DBCPsql.AutoMigrate(&RLoanSector{})
		services.DBCPsql.AutoMigrate(&RLoanDisbursement{})
		services.DBCPsql.AutoMigrate(&RDisbursementHistory{})
		services.DBCPsql.AutoMigrate(&RInstallmentHistory{})
		services.DBCPsql.AutoMigrate(&RInstallmentAccountTransactionCredit{})
		services.DBCPsql.AutoMigrate(&RAreaBranch{})
		services.DBCPsql.AutoMigrate(&RAreaUserMis{})
		services.DBCPsql.AutoMigrate(&RBranchAgent{})
		services.DBCPsql.AutoMigrate(&RBranchUserMis{})
		services.DBCPsql.AutoMigrate(&RGroupAgent{})
		services.DBCPsql.AutoMigrate(&RUserMisRole{})
		services.DBCPsql.AutoMigrate(&RAdjustmentSubmittedBy{})
		services.DBCPsql.AutoMigrate(&RAdjustmentApprovedBy{})
		services.DBCPsql.AutoMigrate(&RInvestorCashout{})
		services.DBCPsql.AutoMigrate(&RCashoutHistory{})
		services.DBCPsql.AutoMigrate(&RUserMisAccessToken{})
		services.DBCPsql.AutoMigrate(&RAgentAccessToken{})
		services.DBCPsql.AutoMigrate(&RGroupBranch{})
		services.DBCPsql.AutoMigrate(&RGroupBorrower{})
		services.DBCPsql.AutoMigrate(&RAccountTransactionDebitInstallment{})
		services.DBCPsql.AutoMigrate(&RAccountTransactionCreditLoan{})
		services.DBCPsql.AutoMigrate(&RInvestorProductPricingLoan{})
		services.DBCPsql.AutoMigrate(&RInvestorProductPricing{})
		services.DBCPsql.AutoMigrate(&RLoanOrder{})
		services.DBCPsql.AutoMigrate(&RAccountTransactionCreditCashout{})
		services.DBCPsql.AutoMigrate(&RAdjustmentAccountTransactionDebit{})
		services.DBCPsql.AutoMigrate(&RInstallmentAdjustment{})
		services.DBCPsql.AutoMigrate(&loanRaw.LoanRaw{})
		services.DBCPsql.AutoMigrate(&RLoanOrderCampaign{})
	} else {
		fmt.Println("Skipp AutoMigrate All Model R Controller")
	}

	services.BaseCrudInitWithDomain("r-cif-access-token", RCifAccessToken{}, []RCifAccessToken{})

	services.BaseCrudInitWithDomain("r-cif-borrower", RCifBorrower{}, []RCifBorrower{})

	services.BaseCrudInitWithDomain("r-cif-investor", RCifInvestor{}, []RCifInvestor{})

	services.BaseCrudInitWithDomain("r-account-transaction-credit", RAccountTransactionCredit{}, []RAccountTransactionCredit{})

	services.BaseCrudInitWithDomain("r-account-transaction-debit", RAccountTransactionDebit{}, []RAccountTransactionDebit{})

	services.BaseCrudInitWithDomain("r-notification-borrower", RNotificationBorrower{}, []RNotificationBorrower{})

	services.BaseCrudInitWithDomain("r-notification-investor", RNotificationInvestor{}, []RNotificationInvestor{})

	services.BaseCrudInitWithDomain("r-account-borrower", RAccountBorrower{}, []RAccountBorrower{})

	services.BaseCrudInitWithDomain("r-account-investor", RAccountInvestor{}, []RAccountInvestor{})

	services.BaseCrudInitWithDomain("r-investor-virtual-account", RInvestorVirtualAccount{}, []RInvestorVirtualAccount{})

	services.BaseCrudInitWithDomain("r-virtual-account-statement", RVirtualAccountStatement{}, []RVirtualAccountStatement{})

	services.BaseCrudInitWithDomain("r-loan-borrower", RLoanBorrower{}, []RLoanBorrower{})

	services.BaseCrudInitWithDomain("r-loan-transaction-credit", RLoanAccountTransactionCredit{}, []RLoanAccountTransactionCredit{})

	services.BaseCrudInitWithDomain("r-loan-history", RLoanHistory{}, []RLoanHistory{})

	services.BaseCrudInitWithDomain("r-loan-area", RLoanArea{}, []RLoanArea{})

	services.BaseCrudInitWithDomain("r-loan-branch", RLoanBranch{}, []RLoanBranch{})

	services.BaseCrudInitWithDomain("r-loan-monitoring", RLoanMonitoring{}, []RLoanMonitoring{})

	services.BaseCrudInitWithDomain("r-loan-group", RLoanGroup{}, []RLoanGroup{})

	services.BaseCrudInitWithDomain("r-loan-installment", RLoanInstallment{}, []RLoanInstallment{})

	services.BaseCrudInitWithDomain("r-loan-sector", RLoanSector{}, []RLoanSector{})

	services.BaseCrudInitWithDomain("r-loan-disbursement", RLoanDisbursement{}, []RLoanDisbursement{})

	services.BaseCrudInitWithDomain("r-disbursement-history", RDisbursementHistory{}, []RDisbursementHistory{})

	services.BaseCrudInitWithDomain("r-installment-history", RInstallmentHistory{}, []RInstallmentHistory{})

	services.BaseCrudInitWithDomain("r-installment-transaction-debit", RInstallmentAccountTransactionCredit{}, []RInstallmentAccountTransactionCredit{})

	services.BaseCrudInitWithDomain("r-area-branch", RAreaBranch{}, []RAreaBranch{})

	services.BaseCrudInitWithDomain("r-area-user-mis", RAreaUserMis{}, []RAreaUserMis{})

	services.BaseCrudInitWithDomain("r-branch-agent", RBranchAgent{}, []RBranchAgent{})

	services.BaseCrudInitWithDomain("r-branch-user-mis", RBranchUserMis{}, []RBranchUserMis{})

	services.BaseCrudInitWithDomain("r-group-agent", RGroupAgent{}, []RGroupAgent{})

	services.BaseCrudInitWithDomain("r-user-mis-role", RUserMisRole{}, []RUserMisRole{})

	services.BaseCrudInitWithDomain("r-adjustment-submitted-by", RAdjustmentSubmittedBy{}, []RAdjustmentSubmittedBy{})

	services.BaseCrudInitWithDomain("r-adjustment-approved-by", RAdjustmentApprovedBy{}, []RAdjustmentApprovedBy{})

	services.BaseCrudInitWithDomain("r-investor-cashout", RInvestorCashout{}, []RInvestorCashout{})

	services.BaseCrudInitWithDomain("r-cashout-history", RCashoutHistory{}, []RCashoutHistory{})

	services.BaseCrudInitWithDomain("r-user-mis-access-token", RUserMisAccessToken{}, []RUserMisAccessToken{})

	services.BaseCrudInitWithDomain("r-agent-access-token", RAgentAccessToken{}, []RAgentAccessToken{})

	services.BaseCrudInitWithDomain("r-group-branch", RGroupBranch{}, []RGroupBranch{})

	services.BaseCrudInitWithDomain("r-group-borrower", RGroupBorrower{}, []RGroupBorrower{})

	services.BaseCrudInitWithDomain("r-account-transaction-debit-installment", RAccountTransactionDebitInstallment{}, []RAccountTransactionDebitInstallment{})

	services.BaseCrudInitWithDomain("r-account-transaction-credit-loan", RAccountTransactionCreditLoan{}, []RAccountTransactionCreditLoan{})

	services.BaseCrudInitWithDomain("r-investor-product-pricing-loan", RInvestorProductPricingLoan{}, []RInvestorProductPricingLoan{})

	services.BaseCrudInitWithDomain("r-investor-product-pricing", RInvestorProductPricing{}, []RInvestorProductPricing{})

	services.BaseCrudInitWithDomain("r-loan-order", RLoanOrder{}, []RLoanOrder{})

	services.BaseCrudInitWithDomain("r-account-transaction-credit-cashout", RAccountTransactionCreditCashout{}, []RAccountTransactionCreditCashout{})

	services.BaseCrudInitWithDomain("r-adjustment-account-transaction-debit", RAdjustmentAccountTransactionDebit{}, []RAdjustmentAccountTransactionDebit{})

	services.BaseCrudInitWithDomain("r-installment-adjustment", RInstallmentAdjustment{}, []RInstallmentAdjustment{})

	services.BaseCrudInitWithDomain("r-loan-raw", loanRaw.LoanRaw{}, []loanRaw.LoanRaw{})

	services.BaseCrudInitWithDomain("r-loan-order-campaign", RLoanOrderCampaign{}, []RLoanOrderCampaign{})
}
