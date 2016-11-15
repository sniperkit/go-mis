package branch

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /branch
func Get(ctx *iris.Context) {
	branch := []Branch{}
	res := services.DBCPsql.Find(&branch)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /branch/get/:id
func GetById(ctx *iris.Context) {
	branch := Branch{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&branch)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /branch/q
func GetByQuery(ctx *iris.Context) {
	branch := []Branch{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&branch)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /branch
func Post(ctx *iris.Context) {
	branch := Branch{}

	err := ctx.ReadJSON(&branch)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&branch)

	ctx.JSON(iris.StatusOK, iris.Map{"data": branch})
}

// PUT /branch/set/:id
func UpdateById(ctx *iris.Context) {
	branch := Branch{}

	err := ctx.ReadJSON(&branch)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&branch).Where("id = ?", id).Update(&branch)

	ctx.JSON(iris.StatusOK, iris.Map{"data": branch})
}

// DELETE /branch/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Branch{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/branch"

	services.DBCPsql.AutoMigrate(&Branch{})

	branch := iris.Party(BASE_URL)
	{
		branch.Get("", Get)
		branch.Get("/get/:id", GetById)
		branch.Get("/q", GetByQuery)
		branch.Post("", Post)
		branch.Put("/set/:id", UpdateById)
		branch.Delete("/delete/:id", DeleteById)
	}
}
