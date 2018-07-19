package sector

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
		services.BaseCrudInit(Sector{}, []Sector{})
}

func GetSectorById(ctx *iris.Context){
	id := ctx.Get("id")
	result := Sector{}
	query := "SELECT sector.\"id\" as \"id\", sector.\"name\" as \"name\", sector.\"description\" as \"description\" FROM sector WHERE sector.\"id\" = ?"
	services.DBCPsql.Raw(query, id).Scan(&result)
	ctx.JSON(iris.StatusOK, iris.Map{"data": result})
}
