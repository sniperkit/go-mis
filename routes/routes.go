package routes

import (
	"fmt"

	"bitbucket.org/go-mis/config"

	"bitbucket.org/go-mis/modules/account"
	"bitbucket.org/go-mis/modules/agent"
	"bitbucket.org/go-mis/modules/borrower"
	"bitbucket.org/go-mis/modules/campaign"
	"bitbucket.org/go-mis/modules/cif"
	"bitbucket.org/go-mis/modules/disbursement"
	"bitbucket.org/go-mis/modules/group"
	"bitbucket.org/go-mis/modules/installment"
	"bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/modules/loan"
	"bitbucket.org/go-mis/modules/notification"
	"bitbucket.org/go-mis/modules/wallet"
	"bitbucket.org/go-mis/modules/wallet-transaction"
)

// If domain is NOT specified,
// then it will automatically initialize all domain
func initializeAll() {
	fmt.Println("[INFO] Initializing all domain")
	account.Init(config.DefaultApiPath)
	agent.Init(config.DefaultApiPath)
	borrower.Init(config.DefaultApiPath)
	campaign.Init(config.DefaultApiPath)
	cif.Init(config.DefaultApiPath)
	disbursement.Init(config.DefaultApiPath)
	group.Init(config.DefaultApiPath)
	installment.Init(config.DefaultApiPath)
	investor.Init(config.DefaultApiPath)
	loan.Init(config.DefaultApiPath)
	notification.Init(config.DefaultApiPath)
	wallet.Init(config.DefaultApiPath)
	walletTransaction.Init(config.DefaultApiPath)
	fmt.Println("[INFO] All domain have been initialized")
}

// Initialize routes
func Init() {
	switch config.Domain {
	case "account":
		account.Init(config.DefaultApiPath)
	case "agent":
		agent.Init(config.DefaultApiPath)
	case "borrower":
		borrower.Init(config.DefaultApiPath)
	case "campaign":
		campaign.Init(config.DefaultApiPath)
	case "cif":
		cif.Init(config.DefaultApiPath)
	case "disbursement":
		disbursement.Init(config.DefaultApiPath)
	case "group":
		group.Init(config.DefaultApiPath)
	case "installment":
		installment.Init(config.DefaultApiPath)
	case "investor":
		investor.Init(config.DefaultApiPath)
	case "loan":
		loan.Init(config.DefaultApiPath)
	case "notification":
		notification.Init(config.DefaultApiPath)
	case "wallet":
		wallet.Init(config.DefaultApiPath)
	case "wallet-transaction":
		walletTransaction.Init(config.DefaultApiPath)
	default:
		initializeAll()
	}
}
