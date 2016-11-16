package account

import (
	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func test(ctx *iris.Context) {
	ctx.Write("haeee")
}

func Init() {
	baseUrl := config.DefaultApiPath + "/" + config.Domain
	services.DBCPsql.AutoMigrate(&Account{})
	services.BaseCrudInit(Account{}, []Account{})

	account := iris.Party(baseUrl)
	{
		account.Get("/test", test)
	}
}
