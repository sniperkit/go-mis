package virtualAccount

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&VirtualAccount{})
	services.BaseCrudInit(VirtualAccount{}, []VirtualAccount{})
}
