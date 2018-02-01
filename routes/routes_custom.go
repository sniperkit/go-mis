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
	"bitbucket.org/go-mis/modules/data-transfer"
	"bitbucket.org/go-mis/modules/disbursement"
	"bitbucket.org/go-mis/modules/disbursement-report"
	"bitbucket.org/go-mis/modules/emergency-loan"
	"bitbucket.org/go-mis/modules/group"
	"bitbucket.org/go-mis/modules/healthycheck"
	"bitbucket.org/go-mis/modules/installment"
	"bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/modules/investor-check"
	"bitbucket.org/go-mis/modules/loan"
	loanInsurance "bitbucket.org/go-mis/modules/loan-insurance"
	"bitbucket.org/go-mis/modules/loan-order"
	"bitbucket.org/go-mis/modules/loan-raw"
	"bitbucket.org/go-mis/modules/location"
	mitramanagement "bitbucket.org/go-mis/modules/mitra-management"
	"bitbucket.org/go-mis/modules/multi-loan"
	"bitbucket.org/go-mis/modules/notification"
	plottingBorrower "bitbucket.org/go-mis/modules/plotting-borrower"
	"bitbucket.org/go-mis/modules/product-pricing"
	prospectiveBorrower "bitbucket.org/go-mis/modules/prospective-borrower"
	"bitbucket.org/go-mis/modules/reports"
	"bitbucket.org/go-mis/modules/sector"
	"bitbucket.org/go-mis/modules/survey"
	systemParameter "bitbucket.org/go-mis/modules/system-parameter"
	"bitbucket.org/go-mis/modules/topsheet"
	"bitbucket.org/go-mis/modules/transaction"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/modules/validation-teller"
	"bitbucket.org/go-mis/modules/virtual-account-statement"
	"bitbucket.org/go-mis/modules/voucher"
	"gopkg.in/iris-contrib/middleware.v4/cors"
	"gopkg.in/kataras/iris.v4"
)

var baseURL = "/api/v2"
var baseRestrictedURL = "/api/restrict"

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
	iris.Any(baseURL+"/installment-approve-success/:stageFrom/:stageTo", installment.SubmitInstallmentByGroupIDAndTransactionDateWithStatus)
	iris.Any(baseURL+"/installment-approve-success-custom-brooooo", installment.SubmitInstallmentByGroupIDAndTransactionDateWithStatusAndInstallmentId)
	//iris.Any(baseURL+"/investor-check/view", investorCheck.FetchDatatables)

	v2 := iris.Party(baseURL, auth.EnsureAuth)
	{
		v2.Any("/me-user-mis", auth.CurrentUserMis)
		v2.Any("/me-notif", auth.Nofif)
		v2.Any("/update-user-branch/:branch_id", userMis.UpdateUserBranch)
		// v2.Any("/me-agent", auth.CurrentAgent)
		v2.Any("/branch", branch.FetchAll)
		v2.Any("/branch-area/:id", branch.IristGetByAreaId)
		v2.Any("/branch/delete/:id", branch.DeleteSingle)
		v2.Any("/branch/detail/:id", branch.GetBranchById)
		v2.Any("/branch/:id", branch.GetByID)
		v2.Any("/branch/area/:id", branch.GetBranchAreaWithoutManager)
		v2.Any("/area", area.FetchAll)
		v2.Any("/area/detail/:id", area.GetByID)
		v2.Any("/area/branch/:id", area.GetByBranch)
		v2.Any("/area/:id", area.GetByIdAreaManager)
		v2.Any("/cif", cif.FetchAll)
		v2.Any("/cif/borrower/:id", cif.GetCifBorrower)
		v2.Any("/cif/investor/:id", cif.GetCifInvestor)
		v2.Any("/investor/set/:investorId/cif/:cifId", cif.UpdateInvestorCif)
		v2.Any("/cif/get/summary", cif.GetCifSummary)
		v2.Any("/group", group.FetchAll)
		v2.Any("/group-by-branch/:branch_id", group.GetGroupByBranchID)
		v2.Any("/group/get/:id", group.GroupDetail)
		v2.Any("/group/create", group.Create)
		v2.Any("/group/set/:id", group.Update)
		v2.Any("/group-borrower/set/:id", group.UpdateGroupBorrower)
		v2.Any("/group/s/branch/:branchId/:searchStr", group.SearchGroup)
		v2.Any("/loan", loan.FetchAll)
		v2.Any("/loan/get/:id", loan.GetLoanDetail)
		v2.Any("/loan/set/:id/stage/:stage", loan.UpdateStage)
		v2.Any("/loan/akad/:id", loan.GetAkadData)
		v2.Any("/loan-stage-history/:id", loan.GetLoanStageHistory)
		v2.Any("/loan-insurance", loanInsurance.GetLoanWithInsurance)
		v2.Any("/loan-insurance/get-finance-report", loanInsurance.GetFinanceReport)
		v2.Any("/loan-insurance/request-refund/:loan_id", loanInsurance.RequestRefund)
		v2.Any("/loan-insurance/apply-refund/:loan_id", loanInsurance.ApplyRefund)
		//v2.Any("/loan-order/pending-waiting", loanOrder.FetchAllPendingWaiting)
		//v2.Any("/loan-order/:orderNo/accept", loanOrder.Accept)
		//v2.Any("/loan-order/:orderNo/reject", loanOrder.Reject)
		v2.Any("/installment", installment.FetchAll)
		v2.Any("/installment-by-type/:type", installment.FetchByType)
		v2.Any("/installment/group/:group_id/by-transaction-date/:transaction_date/stage/:stage", installment.GetInstallmentByGroupIDAndTransactionDate)
		v2.Any("/installment/group/:group_id/by-transaction-date/:transaction_date/submit/:stageFrom/:stageTo", installment.SubmitInstallmentByGroupIDAndTransactionDateWithStatus)
		v2.Any("/installment/submit/:installment_id/status/:status", installment.SubmitInstallmentByInstallmentIDWithStatus)
		v2.Any("/disbursement", disbursement.FetchAll)
		v2.Any("/disbursement/setDate/lwkUpkDate", disbursement.SetLWKUPKDate)
		v2.Any("/disbursement/set/:loan_id/stage/:stage", disbursement.UpdateDisbursementStage)
		v2.Any("/disbursement/get/branch/:branch_id/group/:group_id/disbursement-date/:disbursement_date", disbursement.GetDisbursementDetailByGroup)
		v2.Any("/user-mis", userMis.FetchUserMisAreaBranchRole)
		v2.Any("/user-mis/create", userMis.CreateUserMis)
		v2.Any("/user-mis/get/:id", userMis.GetUserMisById)
		v2.Any("/user-mis/update/:id", userMis.UpdateUserMisById)
		v2.Any("/user-mis/update-password/:id", userMis.UpdateUserMisPasswordById)
		v2.Any("/user-mis/delete/:id", userMis.DeleteUserMis)
		v2.Any("/notification", notification.SendPush)
		v2.Any("/cashout", cashout.FetchAll)
		v2.Any("/cashout/set/:cashout_id/stage/:stage", cashout.UpdateStage)
		v2.Any("/survey", survey.GetProspectiveBorrower)
		v2.Any("/survey/archived", survey.GetProspectiveBorrowerArchived)
		v2.Any("/survey/get/:id", survey.GetProspectiveBorrowerDetail)
		v2.Any("/borrower/approve-survey/:source-type", borrower.Approve)
		v2.Any("/borrower/group/:groupId", borrower.GetBorrowerByGroup)
		v2.Any("/borrower/approve/update-status/:id", borrower.ProspectiveBorrowerUpdateStatus)
		v2.Any("/borrower/reject/update-status/:id", borrower.ProspectiveBorrowerUpdateStatusToReject)
		v2.Get("/borrower/total-by-branch/:branch_id", borrower.GetTotalBorrowerByBranchID)
		v2.Any("/virtual-account-statement", virtualAccountStatement.GetVAStatement)
		v2.Any("/agent", agent.GetAllAgentByBranchID)
		v2.Any("/agent/branch/:id", agent.GetAllAgent)
		v2.Any("/agent/detail/:id", agent.GetAgentById)
		v2.Any("/agent/create", agent.CreateAgent)
		v2.Any("/agent/set/:id", agent.UpdateAgent)
		v2.Any("/agent/update-password/:id", agent.UpdateAgentPasswordById)
		v2.Any("/investor-check/datatables", investorCheck.FetchDatatables)
		v2.Any("/investor-check/validate/:id/status/:status", investorCheck.Validate)
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
		// DAS-276
		v2.Any("/loan-order/accept/:orderNo", loanOrder.AcceptLoanOrder)
		v2.Any("/loan-order/reject/:orderNo", loanOrder.RejectLoanOrder)
		v2.Any("/cif-investor-account", cif.GetCifInvestorAccount)
		v2.Any("/assign-investor-to-loan", loan.AssignInvestorToLoan)
		v2.Any("/product-pricing/s/investor/:searchStr", productPricing.SearchInvestor)
		v2.Any("/product-pricing", productPricing.Create)
		v2.Any("/product-pricing/get/:id", productPricing.GetInvestorsByProductPricing)
		v2.Any("/product-pricing/delete/:id", productPricing.DeleteProductPricing)
		v2.Any("/product-pricing/investor/:invId/productPricing/:ppId/delete", productPricing.DeleteProductPricingInvestor)
		v2.Any("/prospective-borrower", prospectiveBorrower.GetProspectiveBorrower)
		v2.Any("/prospective-borrower/archived", prospectiveBorrower.GetArchivedProspectiveBorrower)
		v2.Any("/prospective-borrower/get/:id", prospectiveBorrower.GetProspectiveBorrowerDetail)
		v2.Any("/prospective-borrower/set/:id/status/:status", prospectiveBorrower.UpdateStatusProspectiveBorrower)
		v2.Any("/sector/detail/:id", sector.GetSectorById)
		v2.Any("/installment-review/get/:branch_id/day/:schedule_day", installment.GetPendingInstallment)
		v2.Any("/installment-review/by-group/:group_id", installment.GetPendingInstallmentDetail)
		v2.Any("/installment-review/set/:installment_id", installment.UpdateInstallmentByInstallmentID)

		v2.Any("/reports/agent", reports.AgentRekap)
		v2.Any("/emergency-loan/borrower/by-branch/:branch_id/available", emergency_loan.FetchAllAvailableBorrower)
		v2.Any("/emergency-loan/submit", emergency_loan.SubmitEmergencyLoan)

		v2.Any("/loan-raw/:id", loanRaw.GetLoanRawById)
		v2.Any("/disbursement-weekly-report", disbursementReport.FetchAllActive)
		v2.Any("/disbursement-weekly-report/:id/detail", disbursementReport.GetDetail)

		v2.Get("/multi-loan-undisbursed", multiloan.GetAllUndisbursedMultiLoan)
		v2.Get("/multi-loan-update-disb-date/:loan_id/:last_disb_date/:next_disb_date", disbursement.UpdateDisbursementDate)
		// Validation Teller
		v2.Any("/validation-teller/save", validationTeller.SaveValidationTeller)
		v2.Any("/validation-teller/getdata", validationTeller.GetDataValidationTeller)
		v2.Any("/validation-teller/detail", validationTeller.GetValidationTellerDetail)
		v2.Any("/validation-teller/detail/save", validationTeller.SaveValidationTellerDetail)
		v2.Any("/installment-pending/get/:currentStage/:branchId/:date", installment.GetPendingInstallmentNew)

		v2.Any("/validation-teller/borrower-notes/save", validationTeller.SaveValidationTellerDetail)
		v2.Any("/reject-notes/:status/:stage/save", validationTeller.SaveRejectNotes)
		v2.Any("/reject-notes/:status/:stage/get/:groupId/:date", validationTeller.GetRejectNotes)

		v2.Any("/validation-teller/group-notes/:logType/save", validationTeller.SaveValidationTellerNotes)
		v2.Any("/validation-teller/view/branch/:branchId/date/:date", validationTeller.GetDataValidationAndTransfer)
		v2.Any("/data-transfer/save", dataTransfer.Save)

		// Mitra Management
		v2.Any("/mitra-management/borrower/:borrowerType/date/:date", mitramanagement.GetBorrowerByInstallmentTypeAndDate)
		v2.Any("/mitra-management/borrower-details", mitramanagement.GetBorrowerDetailByInstallmentTypeAndDate)
		v2.Any("/mitra-management/borrower-status/:statusId/reasons", mitramanagement.GetBorrowerStatusReason)
		v2.Any("/mitra-management/status", mitramanagement.GetStatusAll)
		v2.Any("/mitra-management/submit-reason", mitramanagement.SubmitReason)

		// Plotting Borrower
		v2.Any("/plotting-borrower/plotting-params/save", plottingBorrower.SavePlottingParams)
		v2.Any("/plotting-borrower/plotting-params/all", plottingBorrower.ListPlottingParams)
		v2.Any("/plotting-borrower/plotting-params/find-eligble-investor/:investorId", plottingBorrower.FindEligbleInvestor)
		v2.Any("/plotting-borrower/plotting-params/detail/:investorId", plottingBorrower.GetPlottingParamsDetail)
		v2.Any("/plotting-borrower/plotting-params/toggle-activation", plottingBorrower.TogglePlottingParamsActivation)
		v2.Any("/plotting-borrower/recomended-loan-investor/:investorId", plottingBorrower.FindRecommendedLoanByInvestorCriteria)
		v2.Any("/plotting-borrower/loan-data/:stage", plottingBorrower.FindPlottingBorrower)
		v2.Any("/plotting-borrower/loan-history-scheduler/:date", plottingBorrower.GetSchedulerHistory)
		v2.Any("/plotting-borrower/update-loan-stage", plottingBorrower.UpdateLoanStageHandler)
	}

	vRestrict := iris.Party(baseRestrictedURL, auth.EnsureIp)
	{
		vRestrict.Any("/loan-order/accept/:orderNo", loanOrder.AcceptLoanOrder)
	}

	iris.Get(baseURL+"/generate-topsheet/:group_id", topsheet.GenerateTopsheet)
	iris.Any(baseURL+"/submit-topsheet", topsheet.SubmitTopsheet)
	iris.Any(baseURL+"/submit-topsheet-1", topsheet.SubmitTopsheet)
	iris.Get(baseURL+"/investor-without-va", investor.InvestorWithoutVA)
	iris.Post(baseURL+"/investor-register-va", investor.InvestorRegisterVA)
	iris.Get(baseURL+"/location", location.GetLocation)
	iris.Post(baseURL+"/location/translate", location.TranslateLoc)
	iris.Get(baseURL+"/location/:location_code", location.GetLocationById)
	iris.Get(baseURL+"/healthycheck", healthycheck.Checking)

	iris.Get(baseURL+"/swift", systemParameter.GetSwift)
}
