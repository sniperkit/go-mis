package group

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Group{})
	services.BaseCrudInit(Group{}, []Group{})
}
