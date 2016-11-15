package disbursement

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /disbursement
func Get(ctx *iris.Context) {
	disbursement := []Disbursement{}
	res := services.DBCPsql.Find(&disbursement)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /disbursement/get/:id
func GetById(ctx *iris.Context) {
	disbursement := Disbursement{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&disbursement)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /disbursement/q
func GetByQuery(ctx *iris.Context) {
	disbursement := []Disbursement{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&disbursement)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /disbursement
func Post(ctx *iris.Context) {
	disbursement := Disbursement{}

	err := ctx.ReadJSON(&disbursement)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&disbursement)

	ctx.JSON(iris.StatusOK, iris.Map{"data": disbursement})
}

// PUT /disbursement/set/:id
func UpdateById(ctx *iris.Context) {
	disbursement := Disbursement{}

	err := ctx.ReadJSON(&disbursement)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&disbursement).Where("id = ?", id).Update(&disbursement)

	ctx.JSON(iris.StatusOK, iris.Map{"data": disbursement})
}

// DELETE /disbursement/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Disbursement{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/disbursement"

	services.DBCPsql.AutoMigrate(&Disbursement{})

	disbursement := iris.Party(BASE_URL)
	{
		disbursement.Get("", Get)
		disbursement.Get("/get/:id", GetById)
		disbursement.Get("/q", GetByQuery)
		disbursement.Post("", Post)
		disbursement.Put("/set/:id", UpdateById)
		disbursement.Delete("/delete/:id", DeleteById)
	}
}
