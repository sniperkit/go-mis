package disbursement

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/disbursement"

	services.DomainStructSingle = Disbursement{}
	services.DomainStructArray = []Disbursement{}

	services.DBCPsql.AutoMigrate(&Disbursement{})

	disbursement := iris.Party(BASE_URL)
	{
		disbursement.Get("", services.Get)
		disbursement.Get("/get/:id", services.GetById)
		disbursement.Get("/q", services.GetByQuery)
		disbursement.Post("", services.Post)
		disbursement.Put("/set/:id", services.Put)
		disbursement.Delete("/delete/:id", services.Delete)
	}
}
