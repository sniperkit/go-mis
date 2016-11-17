package routes

import (
	"fmt"

	"bitbucket.org/go-mis/config"

	"bitbucket.org/go-mis/modules/account"
	"bitbucket.org/go-mis/modules/agent"
	"bitbucket.org/go-mis/modules/borrower"
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
	"bitbucket.org/go-mis/modules/product-pricing"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/sector"
	"bitbucket.org/go-mis/modules/wallet"
	"bitbucket.org/go-mis/modules/wallet-transaction"
)

// If domain is NOT specified,
// then it will automatically initialize all domain
func initializeAll() {
	fmt.Println("[INFO] Initializing all domain")
	account.Init()
	agent.Init()
	borrower.Init()
	branch.Init()
	campaign.Init()
	cif.Init()
	disbursement.Init()
	group.Init()
	incentive.Init()
	installment.Init()
	installmentPresence.Init()
	investor.Init()
	loan.Init()
	notification.Init()
	productPricing.Init()
	r.Init()
	sector.Init()
	wallet.Init()
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
	case "product-pricing":
		productPricing.Init()
	case "r":
		r.Init()
	case "sector":
		sector.Init()
	case "wallet":
		wallet.Init()
	case "wallet-transaction":
		walletTransaction.Init()
	default:
		initializeAll()
	}
}
