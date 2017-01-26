package agent

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Agent{})
	services.BaseCrudInit(Agent{}, []Agent{})
}

// FetchAll - fetchAll agent data
func FetchAll(ctx *iris.Context) {
	agents := []Agent{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL").Order("id asc").Find(&agents)
	ctx.JSON(iris.StatusOK, iris.Map{"data": agents})
}

// GetByID agent by id
func GetByID(ctx *iris.Context) {
	agent := Agent{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).First(&agent)
	if agent == (Agent{}) {
		ctx.JSON(iris.StatusOK, iris.Map{"data": iris.Map{}})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": agent})
	}
}
