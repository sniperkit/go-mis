package routes

import (
	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/modules/account"
	"bitbucket.org/go-mis/modules/agent"
	"bitbucket.org/go-mis/modules/area"
	"bitbucket.org/go-mis/modules/auth"
	"bitbucket.org/go-mis/modules/borrower"
	"bitbucket.org/go-mis/modules/branch"
	"bitbucket.org/go-mis/modules/cashout"
	"bitbucket.org/go-mis/modules/cif"
	"bitbucket.org/go-mis/modules/disbursement"
	"bitbucket.org/go-mis/modules/group"
	"bitbucket.org/go-mis/modules/installment"
	"bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/modules/investor-check"
	"bitbucket.org/go-mis/modules/loan"
	"bitbucket.org/go-mis/modules/location"
	"bitbucket.org/go-mis/modules/notification"
	"bitbucket.org/go-mis/modules/survey"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/modules/virtual-account-statement"
	"gopkg.in/iris-contrib/middleware.v4/cors"
	"gopkg.in/kataras/iris.v4"
)

var baseURL = "/api/v2"

// InitCustomApi - initialize custom api
func InitCustomApi() {
	// crs := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{
	// 		"GET", "OPTIONS", "POST",
	// 		"PATCH", "PUT", "DELETE",
	// 	},
	// })
	// app := iris.New()
	// app.Use(crs)
	// Check environment, if `dev` then let the CORS to `*`
	if config.Env == "dev" || config.Env == "development" {
		crs := cors.New(cors.Options{})
		iris.Use(crs)
	}
	iris.Any(baseURL+"/user-mis-login", auth.UserMisLogin)

	v2 := iris.Party(baseURL, auth.EnsureAuth)
	{
		v2.Any("/me-user-mis", auth.CurrentUserMis)
		// v2.Any("/me-agent", auth.CurrentAgent)
		v2.Any("/branch", branch.FetchAll)
		v2.Any("/branch/:id", branch.GetByID)
		v2.Any("/area", area.FetchAll)
		v2.Any("/area/:id", area.GetByID)
		v2.Any("/cif", cif.FetchAll)
		v2.Any("/cif/get/summary", cif.GetCifSummary)
		v2.Any("/group", group.FetchAll)
		v2.Any("/loan", loan.FetchAll)
		v2.Any("/loan/get/:id", loan.GetLoanDetail)
		v2.Any("/loan/set/:id/stage/:stage", loan.UpdateStage)
		v2.Any("/loan/akad/:id", loan.GetAkadData)
		v2.Any("/installment", installment.FetchAll)
		v2.Any("/installment-by-type/:type", installment.FetchByType)
		v2.Any("/installment/group/:group_id/by-transaction-date/:transaction_date", installment.GetInstallmentByGroupIDAndTransactionDate)
		v2.Any("/installment/group/:group_id/by-transaction-date/:transaction_date/submit/:status", installment.SubmitInstallmentByGroupIDAndTransactionDateWithStatus)
		v2.Any("/installment/submit/:installment_id/status/:status", installment.SubmitInstallmentByInstallmentIDWithStatus)
		v2.Any("/disbursement", disbursement.FetchAll)
		v2.Any("/disbursement/set/:loan_id/stage/:stage", disbursement.UpdateDisbursementStage)
		v2.Any("/disbursement/get/branch/:branch_id/group/:group_id/disbursement-date/:disbursement_date", disbursement.GetDisbursementDetailByGroup)
		v2.Any("/user-mis", userMis.FetchUserMisAreaBranchRole)
		v2.Any("/notification", notification.SendPush)
		v2.Any("/cashout", cashout.FetchAll)
		v2.Any("/cashout/set/:cashout_id/stage/:stage", cashout.UpdateStage)
		v2.Any("/survey", survey.GetProspectiveBorrower)
		v2.Any("/survey/get/:id", survey.GetProspectiveBorrowerDetail)
		v2.Any("/borrower/approve", borrower.Approve)
		v2.Any("/borrower/approve/update-status/:id", borrower.ProspectiveBorrowerUpdateStatus)
		v2.Any("/borrower/reject/update-status/:id", borrower.ProspectiveBorrowerUpdateStatusToReject)
		v2.Get("/borrower/total-by-branch/:branch_id", borrower.GetTotalBorrowerByBranchID)
		v2.Any("/virtual-account-statement", virtualAccountStatement.GetVAStatement)
		v2.Any("/agent", agent.GetAllAgentByBranchID)
		v2.Any("/investor-check/datatables", investorCheck.FetchDatatables)
		v2.Any("/investor-check/verify/:id", investorCheck.Verify)
		v2.Get("/investor-for-topup", investor.GetInvestorForTopup)
		v2.Any("/topup/submit", account.DoTopup)
	}

	iris.Get(baseURL+"/investor-without-va", investor.InvestorWithoutVA)
	iris.Post(baseURL+"/investor-register-va", investor.InvestorRegisterVA)
	iris.Get(baseURL+"/location", location.GetLocation)
}
