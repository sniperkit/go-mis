package cif

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /cif
func Get(ctx *iris.Context) {
	cif := []Cif{}
	res := services.DBCPsql.Find(&cif)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /cif/get/:id
func GetById(ctx *iris.Context) {
	cif := Cif{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&cif)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /cif/q
func GetByQuery(ctx *iris.Context) {
	cif := []Cif{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&cif)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /cif
func Post(ctx *iris.Context) {
	cif := Cif{}

	err := ctx.ReadJSON(&cif)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&cif)

	ctx.JSON(iris.StatusOK, iris.Map{"data": cif})
}

// PUT /cif/set/:id
func UpdateById(ctx *iris.Context) {
	cif := Cif{}

	err := ctx.ReadJSON(&cif)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&cif).Where("id = ?", id).Update(&cif)

	ctx.JSON(iris.StatusOK, iris.Map{"data": cif})
}

// DELETE /cif/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Cif{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/cif"

	services.DBCPsql.AutoMigrate(&Cif{})

	cif := iris.Party(BASE_URL)
	{
		cif.Get("", Get)
		cif.Get("/get/:id", GetById)
		cif.Get("/q", GetByQuery)
		cif.Post("", Post)
		cif.Put("/set/:id", UpdateById)
		cif.Delete("/delete/:id", DeleteById)
	}
}
