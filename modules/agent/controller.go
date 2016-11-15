package agent

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /agent
func Get(ctx *iris.Context) {
	agent := []Agent{}
	res := services.DBCPsql.Find(&agent)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /agent/get/:id
func GetById(ctx *iris.Context) {
	agent := Agent{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&agent)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /agent/q
func GetByQuery(ctx *iris.Context) {
	agent := []Agent{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&agent)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /agent
func Post(ctx *iris.Context) {
	agent := Agent{}

	err := ctx.ReadJSON(&agent)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&agent)

	ctx.JSON(iris.StatusOK, iris.Map{"data": agent})
}

// PUT /agent/set/:id
func UpdateById(ctx *iris.Context) {
	agent := Agent{}

	err := ctx.ReadJSON(&agent)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&agent).Where("id = ?", id).Update(&agent)

	ctx.JSON(iris.StatusOK, iris.Map{"data": agent})
}

// DELETE /agent/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Agent{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/agent"

	services.DBCPsql.AutoMigrate(&Agent{})

	agent := iris.Party(BASE_URL)
	{
		agent.Get("", Get)
		agent.Get("/get/:id", GetById)
		agent.Get("/q", GetByQuery)
		agent.Post("", Post)
		agent.Put("/set/:id", UpdateById)
		agent.Delete("/delete/:id", DeleteById)
	}
}
