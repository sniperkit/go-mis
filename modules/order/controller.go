package order

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Order{})
	services.BaseCrudInit(Order{}, []Order{})
}
