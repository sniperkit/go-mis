package account

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/account"

	services.DomainStructSingle = Account{}
	services.DomainStructArray = []Account{}

	services.DBCPsql.AutoMigrate(&Account{})

	account := iris.Party(BASE_URL)
	{
		account.Get("", services.Get)
		account.Get("/get/:id", services.GetById)
		account.Get("/q", services.GetByQuery)
		account.Post("", services.Post)
		account.Put("/set/:id", services.Put)
		account.Delete("/delete/:id", services.Delete)
	}
}
