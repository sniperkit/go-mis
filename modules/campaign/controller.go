package campaign

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/campaign"

	services.DomainStructSingle = Campaign{}
	services.DomainStructArray = []Campaign{}

	services.DBCPsql.AutoMigrate(&Campaign{})

	campaign := iris.Party(BASE_URL)
	{
		campaign.Get("", services.Get)
		campaign.Get("/get/:id", services.GetById)
		campaign.Get("/q", services.GetByQuery)
		campaign.Post("", services.Post)
		campaign.Put("/set/:id", services.Put)
		campaign.Delete("/delete/:id", services.Delete)
	}
}
