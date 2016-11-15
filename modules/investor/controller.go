package investor

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /investor
func Get(ctx *iris.Context) {
	investor := []Investor{}
	res := services.DBCPsql.Find(&investor)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /investor/get/:id
func GetById(ctx *iris.Context) {
	investor := Investor{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&investor)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /investor/q
func GetByQuery(ctx *iris.Context) {
	investor := []Investor{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&investor)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /investor
func Post(ctx *iris.Context) {
	investor := Investor{}

	err := ctx.ReadJSON(&investor)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&investor)

	ctx.JSON(iris.StatusOK, iris.Map{"data": investor})
}

// PUT /investor/set/:id
func UpdateById(ctx *iris.Context) {
	investor := Investor{}

	err := ctx.ReadJSON(&investor)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&investor).Where("id = ?", id).Update(&investor)

	ctx.JSON(iris.StatusOK, iris.Map{"data": investor})
}

// DELETE /investor/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Investor{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/investor"

	services.DBCPsql.AutoMigrate(&Investor{})

	investor := iris.Party(BASE_URL)
	{
		investor.Get("", Get)
		investor.Get("/get/:id", GetById)
		investor.Get("/q", GetByQuery)
		investor.Post("", Post)
		investor.Put("/set/:id", UpdateById)
		investor.Delete("/delete/:id", DeleteById)
	}
}
