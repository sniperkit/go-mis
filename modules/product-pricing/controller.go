package productPricing

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /product-pricing
func Get(ctx *iris.Context) {
	productPricing := []ProductPricing{}
	res := services.DBCPsql.Find(&productPricing)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /product-pricing/get/:id
func GetById(ctx *iris.Context) {
	productPricing := ProductPricing{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&productPricing)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /product-pricing/q
func GetByQuery(ctx *iris.Context) {
	productPricing := []ProductPricing{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&productPricing)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /product-pricing
func Post(ctx *iris.Context) {
	productPricing := ProductPricing{}

	err := ctx.ReadJSON(&productPricing)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&productPricing)

	ctx.JSON(iris.StatusOK, iris.Map{"data": productPricing})
}

// PUT /product-pricing/set/:id
func UpdateById(ctx *iris.Context) {
	productPricing := ProductPricing{}

	err := ctx.ReadJSON(&productPricing)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&productPricing).Where("id = ?", id).Update(&productPricing)

	ctx.JSON(iris.StatusOK, iris.Map{"data": productPricing})
}

// DELETE /product-pricing/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&ProductPricing{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/product-pricing"

	services.DBCPsql.AutoMigrate(&ProductPricing{})

	productPricing := iris.Party(BASE_URL)
	{
		productPricing.Get("", Get)
		productPricing.Get("/get/:id", GetById)
		productPricing.Get("/q", GetByQuery)
		productPricing.Post("", Post)
		productPricing.Put("/set/:id", UpdateById)
		productPricing.Delete("/delete/:id", DeleteById)
	}
}
