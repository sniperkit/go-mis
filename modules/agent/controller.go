package agent

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/agent"

	services.DomainStructSingle = Agent{}
	services.DomainStructArray = []Agent{}

	services.DBCPsql.AutoMigrate(&Agent{})

	agent := iris.Party(BASE_URL)
	{
		agent.Get("", services.Get)
		agent.Get("/get/:id", services.GetById)
		agent.Get("/q", services.GetByQuery)
		agent.Post("", services.Post)
		agent.Put("/set/:id", services.Put)
		agent.Delete("/delete/:id", services.Delete)
	}
}
