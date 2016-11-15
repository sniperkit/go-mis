package walletTransaction

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/wallet-transaction"

	services.DomainStructSingle = WalletTransaction{}
	services.DomainStructArray = []WalletTransaction{}

	services.DBCPsql.AutoMigrate(&WalletTransaction{})

	walletTransaction := iris.Party(BASE_URL)
	{
		walletTransaction.Get("", services.Get)
		walletTransaction.Get("/get/:id", services.GetById)
		walletTransaction.Get("/q", services.GetByQuery)
		walletTransaction.Post("", services.Post)
		walletTransaction.Put("/set/:id", services.Put)
		walletTransaction.Delete("/delete/:id", services.Delete)
	}
}
