package healthycheck

import (
	iris "gopkg.in/kataras/iris.v4"
	"bitbucket.org/go-mis/services"

	"bitbucket.org/go-mis/modules/role"
)

func Checking (ctx *iris.Context){
	query := `SELECT * FROM "sector" LIMIT 1`

	mRole := role.Role{}

	if e := services.DBCPsql.Raw(query).Scan(&mRole).Error; e != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "failed",
			"data":   e,
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   "Database Up",
	})


}