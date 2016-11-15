package installment

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /installment
func Get(ctx *iris.Context) {
	installment := []Installment{}
	res := services.DBCPsql.Find(&installment)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /installment/get/:id
func GetById(ctx *iris.Context) {
	installment := Installment{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&installment)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /installment/q
func GetByQuery(ctx *iris.Context) {
	installment := []Installment{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&installment)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /installment
func Post(ctx *iris.Context) {
	installment := Installment{}

	err := ctx.ReadJSON(&installment)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&installment)

	ctx.JSON(iris.StatusOK, iris.Map{"data": installment})
}

// PUT /installment/set/:id
func UpdateById(ctx *iris.Context) {
	installment := Installment{}

	err := ctx.ReadJSON(&installment)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&installment).Where("id = ?", id).Update(&installment)

	ctx.JSON(iris.StatusOK, iris.Map{"data": installment})
}

// DELETE /installment/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Installment{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/installment"

	services.DBCPsql.AutoMigrate(&Installment{})

	installment := iris.Party(BASE_URL)
	{
		installment.Get("", Get)
		installment.Get("/get/:id", GetById)
		installment.Get("/q", GetByQuery)
		installment.Post("", Post)
		installment.Put("/set/:id", UpdateById)
		installment.Delete("/delete/:id", DeleteById)
	}
}
