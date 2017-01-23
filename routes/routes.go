package routes

import (
	"fmt"

	"bitbucket.org/go-mis/config"

	"bitbucket.org/go-mis/modules/access-token"
	"bitbucket.org/go-mis/modules/account"
	"bitbucket.org/go-mis/modules/account-transaction-credit"
	"bitbucket.org/go-mis/modules/account-transaction-debit"
	"bitbucket.org/go-mis/modules/adjustment"
	"bitbucket.org/go-mis/modules/agent"
	"bitbucket.org/go-mis/modules/area"
	"bitbucket.org/go-mis/modules/borrower"
	"bitbucket.org/go-mis/modules/branch"
	"bitbucket.org/go-mis/modules/cashout"
	"bitbucket.org/go-mis/modules/cashout-history"
	"bitbucket.org/go-mis/modules/cif"
	"bitbucket.org/go-mis/modules/disbursement"
	"bitbucket.org/go-mis/modules/disbursement-history"
	"bitbucket.org/go-mis/modules/group"
	"bitbucket.org/go-mis/modules/installment"
	"bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/modules/loan"
	"bitbucket.org/go-mis/modules/loan-history"
	"bitbucket.org/go-mis/modules/loan-monitoring"
	"bitbucket.org/go-mis/modules/notification"
	"bitbucket.org/go-mis/modules/product-pricing"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/role"
	"bitbucket.org/go-mis/modules/sector"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/modules/virtual-account"
	"bitbucket.org/go-mis/modules/virtual-account-statement"
)

// If domain is NOT specified,
// then it will automatically initialize all domain
func initializeAll() {
	fmt.Println("[INFO] Initializing all domain")

	config.Domain = "access-token"
	accessToken.Init()

	config.Domain = "account"
	account.Init()

	config.Domain = "account-transaction-credit"
	accountTransactionCredit.Init()

	config.Domain = "account-transaction-debit"
	accountTransactionDebit.Init()

	config.Domain = "adjustment"
	adjustment.Init()

	config.Domain = "agent"
	agent.Init()

	config.Domain = "area"
	area.Init()

	config.Domain = "borrower"
	borrower.Init()

	config.Domain = "branch"
	branch.Init()

	config.Domain = "cashout"
	cashout.Init()

	config.Domain = "cashout-history"
	cashoutHistory.Init()

	config.Domain = "cif"
	cif.Init()

	config.Domain = "disbursement"
	disbursement.Init()

	config.Domain = "disbursement-history"
	disbursementHistory.Init()

	config.Domain = "group"
	group.Init()

	config.Domain = "installment"
	installment.Init()

	config.Domain = "investor"
	investor.Init()

	config.Domain = "loan"
	loan.Init()

	config.Domain = "loan-history"
	loanHistory.Init()

	config.Domain = "loan-monitoring"
	loanMonitoring.Init()

	config.Domain = "notification"
	notification.Init()

	config.Domain = "product-pricing"
	productPricing.Init()

	r.Init()

	config.Domain = "role"
	role.Init()

	config.Domain = "sector"
	sector.Init()

	config.Domain = "user-mis"
	userMis.Init()

	config.Domain = "virtual-account"
	virtualAccount.Init()

	config.Domain = "virtual-account-statement"
	virtualAccountStatement.Init()

	fmt.Println("[INFO] All domain have been initialized")
}

// Init - Initialize routes
func Init() {
	switch config.Domain {
	case "access-token":
		accessToken.Init()
	case "account":
		account.Init()
	case "account-transaction-credit":
		accountTransactionCredit.Init()
	case "account-transaction-debit":
		accountTransactionDebit.Init()
	case "adjustment":
		adjustment.Init()
	case "agent":
		agent.Init()
	case "area":
		area.Init()
	case "borrower":
		borrower.Init()
	case "branch":
		branch.Init()
	case "cashout":
		cashout.Init()
	case "cashout-history":
		cashoutHistory.Init()
	case "cif":
		cif.Init()
	case "disbursement":
		disbursement.Init()
	case "disbursement-history":
		disbursementHistory.Init()
	case "group":
		group.Init()
	case "installment":
		installment.Init()
	case "investor":
		investor.Init()
	case "loan":
		loan.Init()
	case "loan-history":
		loanHistory.Init()
	case "loan-monitoring":
		loanMonitoring.Init()
	case "notification":
		notification.Init()
	case "product-pricing":
		productPricing.Init()
	case "r":
		r.Init()
	case "role":
		role.Init()
	case "sector":
		sector.Init()
	case "user-mis":
		userMis.Init()
	case "virtual-account":
		virtualAccount.Init()
	case "virtual-account-statement":
		virtualAccountStatement.Init()
	default:
		initializeAll()
	}
}
