package sector

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/sector"

	services.DomainStructSingle = Sector{}
	services.DomainStructArray = []Sector{}

	services.DBCPsql.AutoMigrate(&Sector{})

	sector := iris.Party(BASE_URL)
	{
		sector.Get("", services.Get)
		sector.Get("/get/:id", services.GetById)
		sector.Get("/q", services.GetByQuery)
		sector.Post("", services.Post)
		sector.Put("/set/:id", services.Put)
		sector.Delete("/delete/:id", services.Delete)
	}
}
