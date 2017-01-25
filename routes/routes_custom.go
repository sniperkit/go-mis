package routes

import (
	"bitbucket.org/go-mis/modules/auth"
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
	}
}
