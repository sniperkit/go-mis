package loan

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/loan"

	services.DomainStructSingle = Loan{}
	services.DomainStructArray = []Loan{}

	services.DBCPsql.AutoMigrate(&Loan{})

	loan := iris.Party(BASE_URL)
	{
		loan.Get("", services.Get)
		loan.Get("/get/:id", services.GetById)
		loan.Get("/q", services.GetByQuery)
		loan.Post("", services.Post)
		loan.Put("/set/:id", services.Put)
		loan.Delete("/delete/:id", services.Delete)
	}
}
