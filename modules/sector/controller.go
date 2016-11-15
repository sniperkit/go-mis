package sector

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /sector
func Get(ctx *iris.Context) {
	sector := []Sector{}
	res := services.DBCPsql.Find(&sector)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /sector/get/:id
func GetById(ctx *iris.Context) {
	sector := Sector{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&sector)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /sector/q
func GetByQuery(ctx *iris.Context) {
	sector := []Sector{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&sector)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /sector
func Post(ctx *iris.Context) {
	sector := Sector{}

	err := ctx.ReadJSON(&sector)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&sector)

	ctx.JSON(iris.StatusOK, iris.Map{"data": sector})
}

// PUT /sector/set/:id
func UpdateById(ctx *iris.Context) {
	sector := Sector{}

	err := ctx.ReadJSON(&sector)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&sector).Where("id = ?", id).Update(&sector)

	ctx.JSON(iris.StatusOK, iris.Map{"data": sector})
}

// DELETE /sector/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Sector{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/sector"

	services.DBCPsql.AutoMigrate(&Sector{})

	sector := iris.Party(BASE_URL)
	{
		sector.Get("", Get)
		sector.Get("/get/:id", GetById)
		sector.Get("/q", GetByQuery)
		sector.Post("", Post)
		sector.Put("/set/:id", UpdateById)
		sector.Delete("/delete/:id", DeleteById)
	}
}
