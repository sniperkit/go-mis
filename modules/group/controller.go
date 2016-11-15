package group

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/group"

	services.DomainStructSingle = Group{}
	services.DomainStructArray = []Group{}

	services.DBCPsql.AutoMigrate(&Group{})

	group := iris.Party(BASE_URL)
	{
		group.Get("", services.Get)
		group.Get("/get/:id", services.GetById)
		group.Get("/q", services.GetByQuery)
		group.Post("", services.Post)
		group.Put("/set/:id", services.Put)
		group.Delete("/delete/:id", services.Delete)
	}
}
