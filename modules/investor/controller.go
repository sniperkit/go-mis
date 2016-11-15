package investor

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/investor"

	services.DomainStructSingle = Investor{}
	services.DomainStructArray = []Investor{}

	services.DBCPsql.AutoMigrate(&Investor{})

	investor := iris.Party(BASE_URL)
	{
		investor.Get("", services.Get)
		investor.Get("/get/:id", services.GetById)
		investor.Get("/q", services.GetByQuery)
		investor.Post("", services.Post)
		investor.Put("/set/:id", services.Put)
		investor.Delete("/delete/:id", services.Delete)
	}
}
