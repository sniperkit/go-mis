package services

import (
	"reflect"
	"strconv"

	"gopkg.in/kataras/iris.v4"
)

var DomainStructSingle interface{}
var DomainStructArray interface{}

// Find all record
func FindAll(m interface{}) {
	sDBCPsql := DBCPsql
	sDBCPsql.Debug().Find(m)
}

// Find by `:id`
func FindById(m interface{}, id string) {
	sDBCPsql := DBCPsql
	sDBCPsql.Where("id = ?", id).Find(m)
}

// Find by query
func FindByQuery(m interface{}, c *iris.Context) {
	sDBCPsql := DBCPsql
	for key, val := range c.URLParams() {
		sDBCPsql = sDBCPsql.Where(key+" LIKE ?", val)
	}
	sDBCPsql.Find(m)
}

// Create new record
func Create(m interface{}) {
	sDBCPsql := DBCPsql
	sDBCPsql.Create(m)
}

// Update record by `:id`
func UpdateById(m interface{}, id string) {
	sDBCPsql := DBCPsql
	sDBCPsql.Model(m).Where("id = ?", id).Update(m)
}

// Delete a record by `:id`
func DeleteById(m interface{}, id int) {
	sDBCPsql := DBCPsql
	sDBCPsql.Where("id = ?", id).Delete(m)
}

// GET /:domain
func Get(ctx *iris.Context) {
	data := reflect.New(reflect.TypeOf(DomainStructArray)).Interface()
	FindAll(data)
	ctx.JSON(iris.StatusOK, iris.Map{"data": &data})
}

// GET /account/get/:id
func GetById(ctx *iris.Context) {
	data := reflect.New(reflect.TypeOf(DomainStructSingle)).Interface()
	FindById(data, ctx.Param("id"))
	ctx.JSON(iris.StatusOK, iris.Map{"data": &data})
}

// GET /:domain/q
func GetByQuery(ctx *iris.Context) {
	data := reflect.New(reflect.TypeOf(DomainStructArray)).Interface()
	FindByQuery(data, ctx)
	ctx.JSON(iris.StatusOK, iris.Map{"data": &data})
}

// POST /:domain
func Post(ctx *iris.Context) {
	data := reflect.New(reflect.TypeOf(DomainStructSingle)).Interface()

	err := ctx.ReadJSON(&data)
	if err != nil {
		panic(err)
	}

	Create(data)
	ctx.JSON(iris.StatusOK, iris.Map{"data": &data})
}

// PUT /:domain/set/:id
func Put(ctx *iris.Context) {
	data := reflect.New(reflect.TypeOf(DomainStructSingle)).Interface()

	err := ctx.ReadJSON(&data)
	if err != nil {
		panic(err)
	}

	UpdateById(data, ctx.Param("id"))

	ctx.JSON(iris.StatusOK, iris.Map{"data": &data})
}

// DELETE /:domain/delete/:id
func Delete(ctx *iris.Context) {
	data := reflect.New(reflect.TypeOf(DomainStructSingle)).Interface()
	id, _ := strconv.Atoi(ctx.Param("id"))
	DeleteById(data, id)
	ctx.JSON(iris.StatusOK, iris.Map{"data": &data})
}
