package routes

import (
	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/modules/account"
	"bitbucket.org/go-mis/modules/adjustment"
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
	"bitbucket.org/go-mis/modules/loan-order"
	"bitbucket.org/go-mis/modules/location"
	"bitbucket.org/go-mis/modules/notification"
	"bitbucket.org/go-mis/modules/survey"
	"bitbucket.org/go-mis/modules/transaction"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/modules/virtual-account-statement"
	"bitbucket.org/go-mis/modules/voucher"
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
	iris.Any(baseURL+"/installment-approve-success/:status", installment.SubmitInstallmentByGroupIDAndTransactionDateWithStatus)

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
		v2.Any("/loan-stage-history/:id", loan.GetLoanStageHistory)
		//v2.Any("/loan-order/pending-waiting", loanOrder.FetchAllPendingWaiting)
		//v2.Any("/loan-order/:orderNo/accept", loanOrder.Accept)
		//v2.Any("/loan-order/:orderNo/reject", loanOrder.Reject)
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
		v2.Any("/survey/archived", survey.GetProspectiveBorrowerArchived)
		v2.Any("/survey/get/:id", survey.GetProspectiveBorrowerDetail)
		v2.Any("/borrower/approve", borrower.Approve)
		v2.Any("/borrower/approve/update-status/:id", borrower.ProspectiveBorrowerUpdateStatus)
		v2.Any("/borrower/reject/update-status/:id", borrower.ProspectiveBorrowerUpdateStatusToReject)
		v2.Get("/borrower/total-by-branch/:branch_id", borrower.GetTotalBorrowerByBranchID)
		v2.Any("/virtual-account-statement", virtualAccountStatement.GetVAStatement)
		v2.Any("/agent", agent.GetAllAgentByBranchID)
		v2.Any("/investor-check/datatables", investorCheck.FetchDatatables)
		v2.Any("/investor-check/verify/:id/status/:status", investorCheck.Verify)
		//v2.Any("/investor-check/verified/:id", investorCheck.Verified)
		v2.Get("/dropping", loan.FetchDropping)
		v2.Any("/dropping/refund/:loan_id/move-stage-to/:stage", loan.RefundAndChangeStageTo)
		v2.Get("/investor-for-topup", investor.GetInvestorForTopup)
		v2.Any("/topup/submit", account.DoTopup)
		v2.Get("/transaction/:type/:investor_id/:start_date/:end_date", transaction.GetData)
		v2.Get("/adjustment", adjustment.GetAdjustment)
		v2.Get("/adjustment/get/:adjustment_id", adjustment.GetAdjustmentDetail)
		v2.Get("/adjustment/installment/:start_date/:end_date", adjustment.GetInReviewInstallment)
		v2.Any("/adjustment/submit/:installment_id", adjustment.SetAdjustmentForInstallment)
		v2.Any("/adjustment/update/:adjustment_id", adjustment.UpdateAdjustmentAndInstallment)
		v2.Any("/submit-adjustment/:account_type", adjustment.SubmitAdjustment)
		v2.Any("/voucher", voucher.FetchAll)
		v2.Any("/loan-order", loanOrder.FetchAll)
		v2.Get("/loan-order/get/:id", loanOrder.FetchSingle)
		v2.Any("/loan-order/accept/:orderNo", loanOrder.AcceptLoanOrder)
		v2.Any("/loan-order/reject/:orderNo", loanOrder.RejectLoanOrder)
		v2.Any("/cif-investor-account", cif.GetCifInvestorAccount)
		v2.Any("/assign-investor-to-loan", loan.AssignInvestorToLoan)
	}

	iris.Get(baseURL+"/investor-without-va", investor.InvestorWithoutVA)
	iris.Post(baseURL+"/investor-register-va", investor.InvestorRegisterVA)
	iris.Get(baseURL+"/location", location.GetLocation)
	iris.Get(baseURL+"/location/:location_code", location.GetLocationById)
}
