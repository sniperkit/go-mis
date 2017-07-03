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
	areaManager := []AreaManager{}

	query := `
	SELECT area."id", area."name", area."city", area."province", 
	"role"."name" as "role",
	coalesce(user_mis."fullname",'') as "manager",r_area_user_mis.*
	FROM area 
	LEFT JOIN r_area_user_mis ON r_area_user_mis."areaId" = area."id" and r_area_user_mis."deletedAt" is null
	LEFT JOIN user_mis ON user_mis."id" = r_area_user_mis."userMisId"
	LEFT JOIN r_user_mis_role ON r_user_mis_role."userMisId" = user_mis."id"
	LEFT JOIN "role" ON "role"."id" = r_user_mis_role."roleId"
	WHERE area."deletedAt" IS NULL AND ("role"."name" LIKE '%Area Manager%' or "role"."name" is NULL) and user_mis."deletedAt" is null
	`

	if e := services.DBCPsql.Raw(query).Find(&areaManager).Error; e != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "failed",
			"data":   e,
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   areaManager,
	})
}

func GetByIdAreaManager(ctx *iris.Context) {
	id := ctx.Get("id")
	area := AreaManager{}

	query := "SELECT area.\"name\", area.\"id\",area.\"city\" AS \"city\", area.\"province\" AS \"province\", user_mis.\"fullname\" AS \"manager\", role.\"name\" AS \"role\" "
	query += "FROM area "
	query += "LEFT JOIN r_area_user_mis ON r_area_user_mis.\"areaId\" = area.\"id\" "
	query += "LEFT JOIN user_mis ON user_mis.\"id\" = r_area_user_mis.\"userMisId\" "
	query += "LEFT JOIN r_user_mis_role ON r_user_mis_role.\"userMisId\" = user_mis.\"id\" "
	query += "LEFT JOIN role ON role.\"id\" = r_user_mis_role.\"roleId\" "
	query += "WHERE area.\"deletedAt\" IS NULL AND area.\"id\" = ? AND (\"role\".\"name\" LIKE '%Area Manager%' OR \"role\".\"name\" IS NULL)"

	if e := services.DBCPsql.Raw(query, id).Scan(&area).Error; e != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "failed",
			"data":   e,
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   area,
	})
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
