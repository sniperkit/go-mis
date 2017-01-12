package accessToken

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&AccessToken{})
	services.BaseCrudInit(AccessToken{}, []AccessToken{})
}
