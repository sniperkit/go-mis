package notification

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GET /notification
func Get(ctx *iris.Context) {
	notification := []Notification{}
	res := services.DBCPsql.Find(&notification)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /notification/get/:id
func GetById(ctx *iris.Context) {
	notification := Notification{}
	res := services.DBCPsql.Where("id = ?", ctx.Param("id")).Find(&notification)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// GET /notification/q
func GetByQuery(ctx *iris.Context) {
	notification := []Notification{}
	res := services.DBCPsql

	for key, val := range ctx.URLParams() {
		res = res.Where(key+" = ?", val)
	}

	res = res.Find(&notification)

	if res.Error != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": res.Error})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": res.Value})
	}
}

// POST /notification
func Post(ctx *iris.Context) {
	notification := Notification{}

	err := ctx.ReadJSON(&notification)
	if err != nil {
		panic(err)
	}

	services.DBCPsql.Create(&notification)

	ctx.JSON(iris.StatusOK, iris.Map{"data": notification})
}

// PUT /notification/set/:id
func UpdateById(ctx *iris.Context) {
	notification := Notification{}

	err := ctx.ReadJSON(&notification)
	if err != nil {
		panic(err)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Model(&notification).Where("id = ?", id).Update(&notification)

	ctx.JSON(iris.StatusOK, iris.Map{"data": notification})
}

// DELETE /notification/delete/:id
func DeleteById(ctx *iris.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}
	services.DBCPsql.Where("id = ?", id).Delete(&Notification{})
}

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/notification"

	services.DBCPsql.AutoMigrate(&Notification{})

	notification := iris.Party(BASE_URL)
	{
		notification.Get("", Get)
		notification.Get("/get/:id", GetById)
		notification.Get("/q", GetByQuery)
		notification.Post("", Post)
		notification.Put("/set/:id", UpdateById)
		notification.Delete("/delete/:id", DeleteById)
	}
}
