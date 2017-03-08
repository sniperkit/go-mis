package account

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Account{})
	services.BaseCrudInit(Account{}, []Account{})
}
