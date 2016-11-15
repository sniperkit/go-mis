package group

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /group
func Get(ctx *iris.Context) {
	group := []Group{}
	res := services.DBCPsql.Find(&group)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /group/get/:id
func GetById(ctx *iris.Context) {
	group := Group{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&group)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /group/q
func GetByQuery(ctx *iris.Context) {
	group := []Group{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&group)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /group
func Post(ctx *iris.Context) {
	group := Group{}

	err := ctx.ReadJSON(&group)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&group)

	ctx.JSON(iris.StatusOK, iris.Map{"data": group})
}

// PUT /group/set/:id
func UpdateById(ctx *iris.Context) {
	group := Group{}

	err := ctx.ReadJSON(&group)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&group).Where("id = ?", id).Update(&group)

	ctx.JSON(iris.StatusOK, iris.Map{"data": group})
}

// DELETE /group/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Group{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/group"

	services.DBCPsql.AutoMigrate(&Group{})

	group := iris.Party(BASE_URL)
	{
		group.Get("", Get)
		group.Get("/get/:id", GetById)
		group.Get("/q", GetByQuery)
		group.Post("", Post)
		group.Put("/set/:id", UpdateById)
		group.Delete("/delete/:id", DeleteById)
	}
}
