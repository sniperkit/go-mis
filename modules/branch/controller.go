package branch

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Branch{})
	services.BaseCrudInit(Branch{}, []Branch{})
}

// FetchAll - fetchAll branchs data
func FetchAll(ctx *iris.Context) {
	bracnhs := []Branch{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL").Order("id asc").Find(&bracnhs)
	ctx.JSON(iris.StatusOK, iris.Map{"data": bracnhs})
}

// GetByID branch bu id
func GetByID(ctx *iris.Context) {
	bracnh := Branch{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).First(&bracnh)
	ctx.JSON(iris.StatusOK, iris.Map{"data": bracnh})

}
