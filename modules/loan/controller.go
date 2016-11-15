package loan

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /loan
func Get(ctx *iris.Context) {
	loan := []Loan{}
	res := services.DBCPsql.Find(&loan)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /loan/get/:id
func GetById(ctx *iris.Context) {
	loan := Loan{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&loan)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /loan/q
func GetByQuery(ctx *iris.Context) {
	loan := []Loan{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&loan)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /loan
func Post(ctx *iris.Context) {
	loan := Loan{}

	err := ctx.ReadJSON(&loan)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&loan)

	ctx.JSON(iris.StatusOK, iris.Map{"data": loan})
}

// PUT /loan/set/:id
func UpdateById(ctx *iris.Context) {
	loan := Loan{}

	err := ctx.ReadJSON(&loan)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&loan).Where("id = ?", id).Update(&loan)

	ctx.JSON(iris.StatusOK, iris.Map{"data": loan})
}

// DELETE /loan/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Loan{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/loan"

	services.DBCPsql.AutoMigrate(&Loan{})

	loan := iris.Party(BASE_URL)
	{
		loan.Get("", Get)
		loan.Get("/get/:id", GetById)
		loan.Get("/q", GetByQuery)
		loan.Post("", Post)
		loan.Put("/set/:id", UpdateById)
		loan.Delete("/delete/:id", DeleteById)
	}
}
