package accountTransactionDebit

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&AccountTransactionDebit{})
	services.BaseCrudInit(AccountTransactionDebit{}, []AccountTransactionDebit{})
}
