package wallet

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/wallet"

	services.DomainStructSingle = Wallet{}
	services.DomainStructArray = []Wallet{}

	services.DBCPsql.AutoMigrate(&Wallet{})

	wallet := iris.Party(BASE_URL)
	{
		wallet.Get("", services.Get)
		wallet.Get("/get/:id", services.GetById)
		wallet.Get("/q", services.GetByQuery)
		wallet.Post("", services.Post)
		wallet.Put("/set/:id", services.Put)
		wallet.Delete("/delete/:id", services.Delete)
	}
}
