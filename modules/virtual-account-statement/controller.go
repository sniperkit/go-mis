package virtualAccountStatement

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&VirtualAccountStatement{})
	services.BaseCrudInit(VirtualAccountStatement{}, []VirtualAccountStatement{})
}
