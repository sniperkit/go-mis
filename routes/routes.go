package routes

import (
	"fmt"

	"bitbucket.org/go-mis/config"

	"bitbucket.org/go-mis/modules/account"
	"bitbucket.org/go-mis/modules/agent"
	"bitbucket.org/go-mis/modules/borrower"
	"bitbucket.org/go-mis/modules/borrower-prospective"
	"bitbucket.org/go-mis/modules/branch"
	"bitbucket.org/go-mis/modules/campaign"
	"bitbucket.org/go-mis/modules/cif"
	"bitbucket.org/go-mis/modules/disbursement"
	"bitbucket.org/go-mis/modules/group"
	"bitbucket.org/go-mis/modules/incentive"
	"bitbucket.org/go-mis/modules/installment"
	"bitbucket.org/go-mis/modules/installment-presence"
	"bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/modules/loan"
	"bitbucket.org/go-mis/modules/notification"
	"bitbucket.org/go-mis/modules/order"
	"bitbucket.org/go-mis/modules/product-pricing"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/role"
	"bitbucket.org/go-mis/modules/sector"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/modules/wallet"
	"bitbucket.org/go-mis/modules/wallet-transaction"
)

// If domain is NOT specified,
// then it will automatically initialize all domain
func initializeAll() {
	fmt.Println("[INFO] Initializing all domain")

	config.Domain = "account"
	account.Init()

	config.Domain = "agent"
	agent.Init()

	config.Domain = "borrower"
	borrower.Init()

	config.Domain = "borrower-prospective"
	borrowerProspective.Init()

	config.Domain = "branch"
	branch.Init()

	config.Domain = "campaign"
	campaign.Init()

	config.Domain = "cif"
	cif.Init()

	config.Domain = "disbursement"
	disbursement.Init()

	config.Domain = "group"
	group.Init()

	config.Domain = "incentive"
	incentive.Init()

	config.Domain = "installment"
	installment.Init()

	config.Domain = "installment-presence"
	installmentPresence.Init()

	config.Domain = "investor"
	investor.Init()

	config.Domain = "loan"
	loan.Init()

	config.Domain = "notification"
	notification.Init()

	config.Domain = "order"
	order.Init()

	config.Domain = "product-pricing"
	productPricing.Init()

	r.Init()

	config.Domain = "role"
	role.Init()

	config.Domain = "sector"
	sector.Init()

	config.Domain = "user-mis"
	userMis.Init()

	config.Domain = "wallet"
	wallet.Init()

	config.Domain = "wallet-transcation"
	walletTransaction.Init()

	fmt.Println("[INFO] All domain have been initialized")
}

// Initialize routes
func Init() {
	switch config.Domain {
	case "account":
		account.Init()
	case "agent":
		agent.Init()
	case "borrower":
		borrower.Init()
	case "borrower-prospective":
		borrowerProspective.Init()
	case "branch":
		branch.Init()
	case "campaign":
		campaign.Init()
	case "cif":
		cif.Init()
	case "disbursement":
		disbursement.Init()
	case "group":
		group.Init()
	case "incentive":
		incentive.Init()
	case "installment":
		installment.Init()
	case "installment-presence":
		installmentPresence.Init()
	case "investor":
		investor.Init()
	case "loan":
		loan.Init()
	case "notification":
		notification.Init()
	case "order":
		order.Init()
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
	case "wallet":
		wallet.Init()
	case "wallet-transaction":
		walletTransaction.Init()
	default:
		initializeAll()
	}
}
