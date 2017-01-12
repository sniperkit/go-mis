package accountTransactionCredit

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&AccountTransactionCredit{})
	services.BaseCrudInit(AccountTransactionCredit{}, []AccountTransactionCredit{})
}
