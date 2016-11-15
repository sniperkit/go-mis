package borrower

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/borrower"

	services.DomainStructSingle = Borrower{}
	services.DomainStructArray = []Borrower{}

	services.DBCPsql.AutoMigrate(&Borrower{})

	borrower := iris.Party(BASE_URL)
	{
		borrower.Get("", services.Get)
		borrower.Get("/get/:id", services.GetById)
		borrower.Get("/q", services.GetByQuery)
		borrower.Post("", services.Post)
		borrower.Put("/set/:id", services.Put)
		borrower.Delete("/delete/:id", services.Delete)
	}
}
