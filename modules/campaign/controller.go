package campaign

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /campaign
func Get(ctx *iris.Context) {
	campaign := []Campaign{}
	res := services.DBCPsql.Find(&campaign)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /campaign/get/:id
func GetById(ctx *iris.Context) {
	campaign := Campaign{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&campaign)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /campaign/q
func GetByQuery(ctx *iris.Context) {
	campaign := []Campaign{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&campaign)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /campaign
func Post(ctx *iris.Context) {
	campaign := Campaign{}

	err := ctx.ReadJSON(&campaign)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&campaign)

	ctx.JSON(iris.StatusOK, iris.Map{"data": campaign})
}

// PUT /campaign/set/:id
func UpdateById(ctx *iris.Context) {
	campaign := Campaign{}

	err := ctx.ReadJSON(&campaign)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&campaign).Where("id = ?", id).Update(&campaign)

	ctx.JSON(iris.StatusOK, iris.Map{"data": campaign})
}

// DELETE /campaign/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Campaign{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/campaign"

	services.DBCPsql.AutoMigrate(&Campaign{})

	campaign := iris.Party(BASE_URL)
	{
		campaign.Get("", Get)
		campaign.Get("/get/:id", GetById)
		campaign.Get("/q", GetByQuery)
		campaign.Post("", Post)
		campaign.Put("/set/:id", UpdateById)
		campaign.Delete("/delete/:id", DeleteById)
	}
}
