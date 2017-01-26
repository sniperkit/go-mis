package area

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Area{})
	services.BaseCrudInit(Area{}, []Area{})
}

// FetchAll - fetchAll agent data
func FetchAll(ctx *iris.Context) {
	areas := []Area{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL").Order("id asc").Find(&areas)
	ctx.JSON(iris.StatusOK, iris.Map{"data": areas})
}

// GetByID agent by id
func GetByID(ctx *iris.Context) {
	area := Area{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).First(&area)
	if area == (Area{}) {
		ctx.JSON(iris.StatusOK, iris.Map{"data": iris.Map{}})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": area})
	}
}
