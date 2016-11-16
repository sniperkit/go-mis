package wallet

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Wallet{})
	services.BaseCrudInit(Wallet{}, []Wallet{})
}
