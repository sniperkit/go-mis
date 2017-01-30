package routes

import (
	"bitbucket.org/go-mis/modules/agent"
	"bitbucket.org/go-mis/modules/area"
	"bitbucket.org/go-mis/modules/auth"
	"bitbucket.org/go-mis/modules/branch"
	"bitbucket.org/go-mis/modules/cif"
	"bitbucket.org/go-mis/modules/group"
	"gopkg.in/kataras/iris.v4"
)

var baseURL = "/api/v2"

// InitCustomApi - initialize custom api
func InitCustomApi() {

	iris.Post(baseURL+"/user-mis-login", auth.UserMisLogin)

	v2 := iris.Party(baseURL, auth.EnsureAuth)
	{
		v2.Get("/me-user-mis", auth.CurrentUserMis)
		v2.Get("/me-agent", auth.CurrentAgent)
		v2.Get("/branchs", branch.FetchAll)
		v2.Get("/branchs/:id", branch.GetByID)
		v2.Get("/agents", agent.FetchAll)
		v2.Get("/agents/:id", agent.GetByID)
		v2.Get("/areas", area.FetchAll)
		v2.Get("/areas/:id", area.GetByID)
		v2.Get("/cifs", cif.FetchAll)
		v2.Get("/cifs/:id", cif.GetByID)
		v2.Get("/groups", group.FetchAll)
	}
}
