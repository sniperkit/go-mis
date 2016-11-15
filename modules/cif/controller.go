package cif

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/cif"

	services.DomainStructSingle = Cif{}
	services.DomainStructArray = []Cif{}

	services.DBCPsql.AutoMigrate(&Cif{})

	cif := iris.Party(BASE_URL)
	{
		cif.Get("", services.Get)
		cif.Get("/get/:id", services.GetById)
		cif.Get("/q", services.GetByQuery)
		cif.Post("", services.Post)
		cif.Put("/set/:id", services.Put)
		cif.Delete("/delete/:id", services.Delete)
	}
}
