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
	areas := []AreaManager{}

	query := "SELECT area.\"id\", area.\"name\", area.\"city\", area.\"province\", user_mis.\"fullname\" as \"managerName\" "
	query += "FROM area "
	query += "JOIN r_area_user_mis ON r_area_user_mis.\"areaId\" = area.\"id\" "
	query += "JOIN user_mis ON user_mis.\"id\" = r_area_user_mis.\"userMisId\" "
	query += "JOIN r_user_mis_role ON r_user_mis_role.\"userMisId\" = user_mis.\"id\" "
	query += "JOIN role ON role.\"id\" = r_user_mis_role.\"roleId\" "
	query += "WHERE role.\"name\" = ?"

	services.DBCPsql.Raw(query, "areamanager").Find(&areas)
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
