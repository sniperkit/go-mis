package auth

import "gopkg.in/kataras/iris.v4"

// Login - login and generate accessToken
func Login(ctx *iris.Context) {

}

// EnsureAuth - validate access token
func EnsureAuth(ctx *iris.Context) {
	ctx.JSON(iris.StatusForbidden, iris.Map{
		"status":  "error",
		"message": "Unauthorized access.",
	})
}

// CurrentUserMis - get current user mis data
func CurrentUserMis(ctx *iris.Context) {

}

// CurrentAgent - get current agent data
func CurrentAgent(ctx *iris.Context) {

}
