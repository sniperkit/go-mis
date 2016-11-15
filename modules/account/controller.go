package account

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /account
func Get(ctx *iris.Context) {
	account := []Account{}
	res := services.DBCPsql.Find(&account)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /account/get/:id
func GetById(ctx *iris.Context) {
	account := Account{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&account)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /account/q
func GetByQuery(ctx *iris.Context) {
	account := []Account{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&account)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /account
func Post(ctx *iris.Context) {
	account := Account{}

	err := ctx.ReadJSON(&account)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&account)

	ctx.JSON(iris.StatusOK, iris.Map{"data": account})
}

// PUT /account/set/:id
func UpdateById(ctx *iris.Context) {
	account := Account{}

	err := ctx.ReadJSON(&account)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&account).Where("id = ?", id).Update(&account)

	ctx.JSON(iris.StatusOK, iris.Map{"data": account})
}

// DELETE /account/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Account{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/account"

	services.DBCPsql.AutoMigrate(&Account{})

	account := iris.Party(BASE_URL)
	{
		account.Get("", Get)
		account.Get("/get/:id", GetById)
		account.Get("/q", GetByQuery)
		account.Post("", Post)
		account.Put("/set/:id", UpdateById)
		account.Delete("/delete/:id", DeleteById)
	}
}
