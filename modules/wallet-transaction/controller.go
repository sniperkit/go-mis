package walletTransaction

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&WalletTransaction{})
	services.BaseCrudInit(WalletTransaction{}, []WalletTransaction{})
}
