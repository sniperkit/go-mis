package routes

import (
	"bitbucket.org/go-mis/modules/area"
	"bitbucket.org/go-mis/modules/auth"
	"bitbucket.org/go-mis/modules/branch"
	"bitbucket.org/go-mis/modules/cif"
	"bitbucket.org/go-mis/modules/disbursement"
	"bitbucket.org/go-mis/modules/group"
	"bitbucket.org/go-mis/modules/installment"
	"bitbucket.org/go-mis/modules/loan"
	"bitbucket.org/go-mis/modules/user-mis"
	"gopkg.in/kataras/iris.v4"
)

var baseURL = "/api/v2"

// InitCustomApi - initialize custom api
func InitCustomApi() {

	iris.Any(baseURL+"/user-mis-login", auth.UserMisLogin)

	v2 := iris.Party(baseURL, auth.EnsureAuth)
	{
		v2.Get("/me-user-mis", auth.CurrentUserMis)
		v2.Get("/me-agent", auth.CurrentAgent)
		v2.Get("/branch", branch.FetchAll)
		v2.Get("/branch/:id", branch.GetByID)
		v2.Get("/area", area.FetchAll)
		v2.Get("/area/:id", area.GetByID)
		v2.Get("/cif", cif.FetchAll)
		v2.Get("/cif/get/summary", cif.GetCifSummary)
		v2.Get("/group", group.FetchAll)
		v2.Get("/loan", loan.FetchAll)
		v2.Get("/loan/get/:id", loan.GetLoanDetail)
		v2.Get("/loan/set/:id/stage/:stage", loan.UpdateStage)
		v2.Get("/installment", installment.FetchAll)
		v2.Post("/installment/submit/:loan_id", installment.SubmitInstallment)
		v2.Get("/disbursement", disbursement.FetchAll)
		v2.Get("/user-mis", userMis.FetchUserMisAreaBranchRole)
	}
}
