package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/modules/access-token"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/role"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func generateAccessToken() string {
	tempTime := fmt.Sprintf("%s", time.Now())
	hashAccessToken := sha256.Sum256([]byte(tempTime))
	accessToken := hex.EncodeToString(hashAccessToken[:])

	return accessToken
}

// UserMisLogin - login for UserMis and generate accessToken
func UserMisLogin(ctx *iris.Context) {
	loginForm := new(LoginForm)

	if err := ctx.ReadJSON(&loginForm); err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": "Bad request.",
		})
		return
	}

	if loginForm.ApiKey != config.ApiKey {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{
			"status":  "error",
			"message": "Unauthorized access.",
		})
		return
	}

	loginForm.HashPassword()
	arrUserMisObj := []userMis.UserMis{}
	services.DBCPsql.Table("user_mis").Where("\"_username\" = ? AND \"_password\" = ?", loginForm.Username, loginForm.Password).Find(&arrUserMisObj)

	if len(arrUserMisObj) == 0 {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{
			"status":  "error",
			"message": "Invalid username/password. Please try again.",
		})
		return
	}

	accessTokenHash := generateAccessToken()
	accessTokenObj := accessToken.AccessToken{Type: loginForm.Type, AccessToken: accessTokenHash}
	services.DBCPsql.Table("access_token").Create(&accessTokenObj)

	userMisObj := arrUserMisObj[0]
	rUserMisAccessToken := r.RUserMisAccessToken{UserMisId: userMisObj.ID, AccessTokenId: accessTokenObj.ID}
	services.DBCPsql.Table("r_user_mis_access_token").Create(&rUserMisAccessToken)

	roleObj := role.Role{}
	queryRole := "SELECT role.* FROM role JOIN r_user_mis_role ON r_user_mis_role.\"roleId\" = role.\"id\" JOIN user_mis ON user_mis.\"id\" = r_user_mis_role.\"userMisId\" WHERE user_mis.\"id\" = ?"
	services.DBCPsql.Raw(queryRole, userMisObj.ID).First(&roleObj)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":      "success",
		"name":        userMisObj.Fullname,
		"accessToken": accessTokenHash,
		"role": iris.Map{
			"assignedRole": roleObj.Name,
			"config":       roleObj.Config,
		},
	})
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
