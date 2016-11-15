package walletTransaction

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /wallet-transaction
func Get(ctx *iris.Context) {
	walletTransaction := []WalletTransaction{}
	res := services.DBCPsql.Find(&walletTransaction)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /wallet-transaction/get/:id
func GetById(ctx *iris.Context) {
	walletTransaction := WalletTransaction{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&walletTransaction)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /wallet-transaction/q
func GetByQuery(ctx *iris.Context) {
	walletTransaction := []WalletTransaction{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&walletTransaction)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /wallet-transaction
func Post(ctx *iris.Context) {
	walletTransaction := WalletTransaction{}

	err := ctx.ReadJSON(&walletTransaction)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&walletTransaction)

	ctx.JSON(iris.StatusOK, iris.Map{"data": walletTransaction})
}

// PUT /wallet-transaction/set/:id
func UpdateById(ctx *iris.Context) {
	walletTransaction := WalletTransaction{}

	err := ctx.ReadJSON(&walletTransaction)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&walletTransaction).Where("id = ?", id).Update(&walletTransaction)

	ctx.JSON(iris.StatusOK, iris.Map{"data": walletTransaction})
}

// DELETE /wallet-transaction/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&WalletTransaction{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/wallet-transaction"

	services.DBCPsql.AutoMigrate(&WalletTransaction{})

	walletTransaction := iris.Party(BASE_URL)
	{
		walletTransaction.Get("", Get)
		walletTransaction.Get("/get/:id", GetById)
		walletTransaction.Get("/q", GetByQuery)
		walletTransaction.Post("", Post)
		walletTransaction.Put("/set/:id", UpdateById)
		walletTransaction.Delete("/delete/:id", DeleteById)
	}
}
