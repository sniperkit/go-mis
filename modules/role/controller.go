package role

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Role{})
	services.BaseCrudInit(Role{}, []Role{})
}
