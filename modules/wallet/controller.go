package wallet

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /wallet
func Get(ctx *iris.Context) {
	wallet := []Wallet{}
	res := services.DBCPsql.Find(&wallet)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /wallet/get/:id
func GetById(ctx *iris.Context) {
	wallet := Wallet{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&wallet)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /wallet/q
func GetByQuery(ctx *iris.Context) {
	wallet := []Wallet{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&wallet)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /wallet
func Post(ctx *iris.Context) {
	wallet := Wallet{}

	err := ctx.ReadJSON(&wallet)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&wallet)

	ctx.JSON(iris.StatusOK, iris.Map{"data": wallet})
}

// PUT /wallet/set/:id
func UpdateById(ctx *iris.Context) {
	wallet := Wallet{}

	err := ctx.ReadJSON(&wallet)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&wallet).Where("id = ?", id).Update(&wallet)

	ctx.JSON(iris.StatusOK, iris.Map{"data": wallet})
}

// DELETE /wallet/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Wallet{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/wallet"

	services.DBCPsql.AutoMigrate(&Wallet{})

	wallet := iris.Party(BASE_URL)
	{
		wallet.Get("", Get)
		wallet.Get("/get/:id", GetById)
		wallet.Get("/q", GetByQuery)
		wallet.Post("", Post)
		wallet.Put("/set/:id", UpdateById)
		wallet.Delete("/delete/:id", DeleteById)
	}
}
