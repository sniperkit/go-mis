package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"regexp"

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
	services.DBCPsql.Table("user_mis").Where("\"_username\" = ? AND \"_password\" = ? AND \"deletedAt\" IS NULL AND (\"isSuspended\" = FALSE OR \"isSuspended\" IS NULL)", loginForm.Username, loginForm.Password).Find(&arrUserMisObj)

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
	queryRole := "SELECT role.* FROM \"role\" JOIN r_user_mis_role ON r_user_mis_role.\"roleId\" = \"role\".\"id\" JOIN user_mis ON user_mis.\"id\" = r_user_mis_role.\"userMisId\" WHERE user_mis.\"id\" = ?"

	services.DBCPsql.Raw(queryRole, userMisObj.ID).Scan(&roleObj)

	if roleObj.ID == 0 {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{
			"status":  "error",
			"message": "Your account doesn't have any role. Please ask your superadmin to assign a role.",
		})
		return
	}

	// for dashboard
	re := regexp.MustCompile("(?i)area\\s*manager") // area manager, Area Manager, ArEaManager are valid
	if re.FindString(roleObj.Name) != "" { // area manager
		// get the area of this user 
		rAreaUserMis := r.RAreaUserMis{} 
		query := `select "areaId" from r_area_user_mis where "userMisId" = ?`
		services.DBCPsql.Raw(query, userMisObj.ID).Scan(&rAreaUserMis)
		
		// get all branches in this area
		type branchType struct {
			Id uint64 `json:"id"`
			Name string `json:"name"`
		}

		branches := []branchType{}
		query = `select branch.id, branch."name" from r_area_branch 
						join branch on branch.Id = r_area_branch.id
						where "areaId" = ?`
		services.DBCPsql.Raw(query, rAreaUserMis.AreaId).Scan(&branches)

		rUserMisBranch := r.RBranchUserMis{}
		services.DBCPsql.Table("r_branch_user_mis").Where(" \"userMisId\" = ? ", userMisObj.ID).First(&rUserMisBranch)

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data": iris.Map{
				"name":        userMisObj.Fullname,
				"accessToken": accessTokenHash,
				"branchId":    rUserMisBranch.BranchId,
				"role": iris.Map{
					"assignedRole": roleObj.Name,
					"config":       roleObj.Config,
				},
				"branches":branches,
				"roleId": roleObj.ID,
				"areaId": rAreaUserMis.AreaId,
			},
		})
	} else {
		rUserMisBranch := r.RBranchUserMis{}
		services.DBCPsql.Table("r_branch_user_mis").Where(" \"userMisId\" = ? ", userMisObj.ID).First(&rUserMisBranch)

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data": iris.Map{
				"name":        userMisObj.Fullname,
				"accessToken": accessTokenHash,
				"branchId":    rUserMisBranch.BranchId,
				"role": iris.Map{
					"assignedRole": roleObj.Name,
					"config":       roleObj.Config,
				},
			},
		})
	}

	
}

// EnsureAuth - validate access token
func EnsureAuth(ctx *iris.Context) {
	accessToken := ctx.URLParam("accessToken")

	userObj := userMis.UserMis{}
	queryAccessToken := "SELECT user_mis.* FROM access_token JOIN r_user_mis_access_token ON r_user_mis_access_token.\"accessTokenId\" = access_token.\"id\" JOIN user_mis ON user_mis.\"id\" = r_user_mis_access_token.\"userMisId\" WHERE access_token.\"accessToken\" = ? "
	queryAccessToken += "AND access_token.\"deletedAt\" IS NULL"
	services.DBCPsql.Raw(queryAccessToken, accessToken).Scan(&userObj)

	if userObj == (userMis.UserMis{}) {
		ctx.JSON(iris.StatusForbidden, iris.Map{
			"status":  "error",
			"message": "Unauthorized access.",
		})
	} else {
		ctx.Set("USER_MIS", userObj)

		rBranchUserMisSchema := r.RBranchUserMis{}

		queryGetBranch := "SELECT r_branch_user_mis.\"branchId\" "
		queryGetBranch += "FROM user_mis "
		queryGetBranch += "JOIN r_branch_user_mis ON r_branch_user_mis.\"userMisId\" = user_mis.id WHERE user_mis.id = ?"

		services.DBCPsql.Raw(queryGetBranch, userObj.ID).Scan(&rBranchUserMisSchema)

		ctx.Set("BRANCH_ID", rBranchUserMisSchema.BranchId)

		ctx.Next()
	}
}

// CurrentUserMis - get current user mis data
func CurrentUserMis(ctx *iris.Context) {
	userMisObj := ctx.Get("USER_MIS").(userMis.UserMis)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"id":   userMisObj.ID,
			"name": userMisObj.Fullname,
		},
	})
}
