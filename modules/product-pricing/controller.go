package productPricing

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&ProductPricing{})
	services.BaseCrudInit(ProductPricing{}, []ProductPricing{})
}
