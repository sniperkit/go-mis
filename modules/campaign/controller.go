package campaign

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Campaign{})
	services.BaseCrudInit(Campaign{}, []Campaign{})
}
