package services

import (
	"reflect"

	"bitbucket.org/go-mis/config"
	"gopkg.in/kataras/iris.v4"
)

type Container struct {
	SingleObj interface{}
	ArrayObj  interface{}
}

// GET /:domain
func Get(model interface{}) func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).ArrayObj)).Interface()
		DBCPsql.Find(m)
		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}

// GET /:domain/get/:id
func GetById(model interface{}) func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).SingleObj)).Interface()
		DBCPsql.Where("id = ?", ctx.Param("id")).Find(m)
		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}

// GET /:domain/q
func GetByQuery(model interface{}) func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).SingleObj)).Interface()

		con := DBCPsql
		for key, val := range ctx.URLParams() {
			con = con.Where(key+" LIKE ?", val)
		}

		con.Find(m)

		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}

// POST /:domain
func Post(model interface{}) func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).SingleObj)).Interface()

		err := ctx.ReadJSON(&m)
		if err != nil {
			panic(err)
		}

		DBCPsql.Create(m)

		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}

// PUT /:domain/set/:id
func Put(model interface{}) func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).SingleObj)).Interface()

		err := ctx.ReadJSON(&m)
		if err != nil {
			panic(err)
		}

		DBCPsql.Model(m).Where("id = ?", ctx.Param("id")).Update(m)
		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}

// DELETE /:domain/delete/:id
func DeleteById(model interface{}) func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).SingleObj)).Interface()
		DBCPsql.Where("id = ?", ctx.Param("id")).Delete(m)
		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}

// Initialize Base CRUD
func BaseCrudInit(singleObj interface{}, arrayObj interface{}) {
	BASE_URL := config.DefaultApiPath + "/" + config.Domain

	model := new(Container)
	model.SingleObj = singleObj
	model.ArrayObj = arrayObj

	crudParty := iris.Party(BASE_URL)
	{
		crudParty.Get("", Get(model))
		crudParty.Get("/get/:id", GetById(model))
		crudParty.Get("/q", GetByQuery(model))
		crudParty.Post("", Post(model))
		crudParty.Put("/set/:id", Put(model))
		crudParty.Delete("/delete/:id", DeleteById(model))
	}
}
