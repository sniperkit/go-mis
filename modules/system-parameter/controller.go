package systemParameter

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&SystemParameter{})
	services.BaseCrudInit(SystemParameter{}, []SystemParameter{})
}
