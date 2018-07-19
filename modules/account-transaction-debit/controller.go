package accountTransactionDebit

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(AccountTransactionDebit{}, []AccountTransactionDebit{})
}
