package borrower

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /borrower
func Get(ctx *iris.Context) {
	borrower := []Borrower{}
	res := services.DBCPsql.Find(&borrower)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /borrower/get/:id
func GetById(ctx *iris.Context) {
	borrower := Borrower{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&borrower)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /borrower/q
func GetByQuery(ctx *iris.Context) {
	borrower := []Borrower{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&borrower)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /borrower
func Post(ctx *iris.Context) {
	borrower := Borrower{}

	err := ctx.ReadJSON(&borrower)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&borrower)

	ctx.JSON(iris.StatusOK, iris.Map{"data": borrower})
}

// PUT /borrower/set/:id
func UpdateById(ctx *iris.Context) {
	borrower := Borrower{}

	err := ctx.ReadJSON(&borrower)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&borrower).Where("id = ?", id).Update(&borrower)

	ctx.JSON(iris.StatusOK, iris.Map{"data": borrower})
}

// DELETE /borrower/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Borrower{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/borrower"

	services.DBCPsql.AutoMigrate(&Borrower{})

	borrower := iris.Party(BASE_URL)
	{
		borrower.Get("", Get)
		borrower.Get("/get/:id", GetById)
		borrower.Get("/q", GetByQuery)
		borrower.Post("", Post)
		borrower.Put("/set/:id", UpdateById)
		borrower.Delete("/delete/:id", DeleteById)
	}
}
