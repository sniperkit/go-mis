package healthycheck

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"

	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/modules/role"
)

func Checking(ctx *iris.Context) {
	query := `SELECT * FROM "sector" LIMIT 1`

	mRole := role.Role{}

	// DB status
	db_stat := make(map[string]string)
	db_stat["dbms"] = "Postgresql"
	if e := services.DBCPsql.Raw(query).Scan(&mRole).Error; e != nil {
		db_stat["status"] = "down"
	}
	db_stat["status"] = "up"

	// API version
	v := config.Version

	// add all required status here, for now it's only the DB
	data := []map[string]string{db_stat}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":     "success",
		"data":       data,
		"APIVersion": v,
	})
}
