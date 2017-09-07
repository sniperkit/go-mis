package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/role"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/services"
	"github.com/dgrijalva/jwt-go"
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

	type Cas struct {
		Username string `json:"username"`
		Password string `json:"password"`
		UserType string `json:"userType"`
	}

	loginForm := new(LoginForm)

	// login := ctx.ReadJSON(&loginForm)
	// fmt.Println(login)
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
	//
	// loginForm.HashPassword()

	u := Cas{Username: loginForm.Username, Password: loginForm.Password, UserType: "MIS"}
	fmt.Println(u)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	fmt.Println(config.GoCasApiPath)
	fmt.Println(config.SignStringKey)
	res, _ := http.Post(config.GoCasApiPath+"/api/v1/auth", "application/json; charset=utf-8", b)

	var casResp struct {
		Status uint64 `json:"status"`
		Data   string `json:"data"`
	}

	json.NewDecoder(res.Body).Decode(&casResp)

	fmt.Println(casResp.Data)

	arrUserMisObj := []userMis.UserMis{}
	services.DBCPsql.Table("user_mis").Where("\"_username\" = ? AND \"deletedAt\" IS NULL AND (\"isSuspended\" = FALSE OR \"isSuspended\" IS NULL)", loginForm.Username).Find(&arrUserMisObj)

	if casResp.Data == "" {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{
			"status":  "error",
			"message": "Invalid username/password. Please try again.",
		})
		return
	}

	accessTokenHash := casResp.Data

	userMisObj := arrUserMisObj[0]

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

	// get the area of this user
	rAreaUserMis := r.RAreaUserMis{}
	query := `select "areaId" from r_area_user_mis where "userMisId" = ?`
	services.DBCPsql.Raw(query, userMisObj.ID).Scan(&rAreaUserMis)

	// for dashboard
	re := regexp.MustCompile("(?i)area\\s*manager") // area manager, Area Manager, ArEaManager are valid
	if re.FindString(roleObj.Name) != "" {          // area manager

		// get all branches in this area
		type branchType struct {
			Id   uint64 `json:"id"`
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
				"idUserMis":   userMisObj.ID,
				"name":        userMisObj.Fullname,
				"accessToken": accessTokenHash,
				"branchId":    rUserMisBranch.BranchId,
				"role": iris.Map{
					"assignedRole": roleObj.Name,
					"config":       roleObj.Config,
				},
				"branches": branches,
				"roleId":   roleObj.ID,
				"areaId":   rAreaUserMis.AreaId,
				// =======
				// rUserMisBranch := r.RBranchUserMis{}
				// services.DBCPsql.Table("r_branch_user_mis").Where(" \"userMisId\" = ? ", userMisObj.ID).First(&rUserMisBranch)

				// rAreaUserMis := r.RAreaUserMis{}
				// services.DBCPsql.Table("r_area_user_mis").Where(" \"userMisId\" = ? ", userMisObj.ID).First(&rAreaUserMis)

				// ctx.JSON(iris.StatusOK, iris.Map{
				// 	"status": "success",
				// 	"data": iris.Map{
				// 		"name":        userMisObj.Fullname,
				// 		"accessToken": accessTokenHash,
				// 		"branchId":    rUserMisBranch.BranchId,
				// 		"areaId":    rAreaUserMis.AreaId,
				// 		"role": iris.Map{
				// 			"assignedRole": roleObj.Name,
				// 			"config":       roleObj.Config,
				// >>>>>>> master
			},
		})
	} else {
		rUserMisBranch := r.RBranchUserMis{}
		services.DBCPsql.Table("r_branch_user_mis").Where(" \"userMisId\" = ? ", userMisObj.ID).First(&rUserMisBranch)

		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data": iris.Map{
				"idUserMis":   userMisObj.ID,
				"username":    userMisObj.Username,
				"name":        userMisObj.Fullname,
				"accessToken": accessTokenHash,
				"branchId":    rUserMisBranch.BranchId,
				"areaId":      rAreaUserMis.AreaId,
				"role": iris.Map{
					"assignedRole": roleObj.Name,
					"config":       roleObj.Config,
				},
			},
		})
	}

}

type CasResponse struct {
	Status string `json:"username"`
	Data   string `json:"username"`
}

// EnsureAuth - validate access token
func EnsureAuth(ctx *iris.Context) {
	accessToken := ctx.URLParam("accessToken")
	if accessToken == "" {
		accessToken = ctx.RequestHeader("accessToken")
	}
	signString := []byte(config.SignStringKey)

	claim := &jwt.StandardClaims{}
	jwt.ParseWithClaims(accessToken, claim, func(token *jwt.Token) (interface{}, error) {
		return signString, nil
	})
	if claim.Id == "" || time.Now().UnixNano() > claim.NotBefore {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{
			"status":  "error",
			"message": "Unauthorized access.",
		})
		return
	}
	userObj := userMis.UserMis{}
	queryAccessToken := "SELECT * FROM user_mis WHERE \"id\" = ? "
	queryAccessToken += "AND \"deletedAt\" IS NULL"
	services.DBCPsql.Raw(queryAccessToken, claim.Id).Scan(&userObj)

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

type Count struct {
	Total int64 `gorm:"column:count" json:"total" `
}

// Nofif - get notif user
func Nofif(ctx *iris.Context) {
	query := `select count(*) from product_pricing where product_pricing."isInstitutional" = false and current_date::date between product_pricing."startDate"::date and product_pricing."endDate"::date`

	countSchema := Count{}
	services.DBCPsql.Raw(query).Scan(&countSchema)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"ppretail_active": countSchema.Total,
		},
	})
}
