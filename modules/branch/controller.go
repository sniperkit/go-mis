package branch

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/branch"

	services.DomainStructSingle = Branch{}
	services.DomainStructArray = []Branch{}

	services.DBCPsql.AutoMigrate(&Branch{})

	branch := iris.Party(BASE_URL)
	{
		branch.Get("", services.Get)
		branch.Get("/get/:id", services.GetById)
		branch.Get("/q", services.GetByQuery)
		branch.Post("", services.Post)
		branch.Put("/set/:id", services.Put)
		branch.Delete("/delete/:id", services.Delete)
	}
}
