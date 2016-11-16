package branch

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Branch{})
	services.BaseCrudInit(Branch{}, []Branch{})
}
