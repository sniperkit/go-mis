package services

import (
	"reflect"
	"time"
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
		DBCPsql.Where("\"deletedAt\" IS NULL").Order("id asc").Find(m)
		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}

// GET /:domain/get/:id
func GetById(model interface{}) func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).SingleObj)).Interface()
		if dbc:= DBCPsql.Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).Find(m); dbc.Error != nil {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": dbc.Error})
			return
		}
		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}


// GET /:domain/q
func GetByQuery(model interface{}) func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).ArrayObj)).Interface()

		con := DBCPsql
		for key, val := range ctx.URLParams() {
			if key != "apiKey" && key != "q" && ctx.URLParam("q") == "like" {
				con = con.Where("\""+key+"\" LIKE ?", val)
			} else if key != "apiKey" && key != "q" && ctx.URLParam("q") == "equal" {
				con = con.Where("\""+key+"\" = ?", val)
			}
		}

		con.Where("\"deletedAt\" IS NULL").Order("id asc").Find(m)

		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}

// count object
func GetSingle(model interface{})  func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).ArrayObj)).Interface()
		con := DBCPsql
		for key, val := range ctx.URLParams() {
			if key != "apiKey" && key != "q" && ctx.URLParam("q") == "like" {
				con = con.Where("\""+key+"\" LIKE ?", val)
			} else if key != "apiKey" && key != "q" && ctx.URLParam("q") == "equal" {
				con = con.Where("\""+key+"\" = ?", val)
			}
		}

		con.Where("\"deletedAt\" IS NULL").Order("id asc limit 1").Find(m)

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

		if dbc := DBCPsql.Create(m); dbc.Error != nil {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": dbc.Error})
			return
		}

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

		DBCPsql.Model(m).Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).Update(m)
		_, isFieldExist := reflect.TypeOf((model.(*Container)).SingleObj).FieldByName("Password")
		if isFieldExist {
			DBCPsql.Model(m).Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).Update(m)
		}

		GetById(model)(ctx)
	}
}

// DELETE /:domain/delete/:id
func DeleteById(model interface{}) func(ctx *iris.Context) {
	return func(ctx *iris.Context) {
		m := reflect.New(reflect.TypeOf((model.(*Container)).SingleObj)).Interface()
		DBCPsql.Model(m).Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).UpdateColumn("deletedAt", time.Now())
		ctx.JSON(iris.StatusOK, iris.Map{"data": m})
	}
}

// Chech Authentication
func CheckAuth(ctx *iris.Context) {
	if ctx.URLParam("apiKey") != config.ApiKey {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{"error": "Unauthorized access."})
		return
	}

	ctx.Next()
}

func CheckAuthForm(ctx *iris.Context) {
	apiKey := ctx.FormValueString("apiKey")

	if apiKey != config.ApiKey {
		ctx.JSON(iris.StatusUnauthorized, iris.Map{"error": "Unauthorized access."})
		return
	}

	ctx.Next()
}

// Initialize Base CRUD
func BaseCrudInit(singleObj interface{}, arrayObj interface{}) {
	BASE_URL := config.DefaultApiPath + "/" + config.Domain

	model := new(Container)
	model.SingleObj = singleObj
	model.ArrayObj = arrayObj

	crudParty := iris.Party(BASE_URL)
	{
		crudParty.Get("", CheckAuth, Get(model))
		crudParty.Get("/get/:id", CheckAuth, GetById(model))
		crudParty.Get("/q", CheckAuth, GetByQuery(model))
		crudParty.Post("", CheckAuthForm, Post(model))
		crudParty.Any("/set/:id", CheckAuth, Put(model))
		crudParty.Any("/delete/:id", CheckAuth, DeleteById(model))
	}
}

// Initialize Base CRUD
func BaseCrudInitWithDomain(domain string, singleObj interface{}, arrayObj interface{}) {
	BASE_URL := config.DefaultApiPath + "/" + domain

	model := new(Container)
	model.SingleObj = singleObj
	model.ArrayObj = arrayObj

	crudParty := iris.Party(BASE_URL)
	{
		crudParty.Get("", CheckAuth, Get(model))
		crudParty.Get("/get/:id", CheckAuth, GetById(model))
		crudParty.Get("/q", CheckAuth, GetByQuery(model))
		crudParty.Post("", CheckAuthForm, Post(model))
		crudParty.Any("/set/:id", CheckAuth, Put(model))
		crudParty.Any("/delete/:id", CheckAuth, DeleteById(model))
		crudParty.Any("/single", CheckAuth, GetSingle(model))
	}
}
