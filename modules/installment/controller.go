package installment

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/installment"

	services.DomainStructSingle = Installment{}
	services.DomainStructArray = []Installment{}

	services.DBCPsql.AutoMigrate(&Installment{})

	installment := iris.Party(BASE_URL)
	{
		installment.Get("", services.Get)
		installment.Get("/get/:id", services.GetById)
		installment.Get("/q", services.GetByQuery)
		installment.Post("", services.Post)
		installment.Put("/set/:id", services.Put)
		installment.Delete("/delete/:id", services.Delete)
	}
}
