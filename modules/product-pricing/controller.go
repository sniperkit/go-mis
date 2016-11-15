package productPricing

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init(defaultApiPath string) {
	BASE_URL := defaultApiPath + "/product-pricing"

	services.DomainStructSingle = ProductPricing{}
	services.DomainStructArray = []ProductPricing{}

	services.DBCPsql.AutoMigrate(&ProductPricing{})

	productPricing := iris.Party(BASE_URL)
	{
		productPricing.Get("", services.Get)
		productPricing.Get("/get/:id", services.GetById)
		productPricing.Get("/q", services.GetByQuery)
		productPricing.Post("", services.Post)
		productPricing.Put("/set/:id", services.Put)
		productPricing.Delete("/delete/:id", services.Delete)
	}
}
