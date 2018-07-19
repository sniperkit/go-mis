package accountTransactionCredit

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(AccountTransactionCredit{}, []AccountTransactionCredit{})
}
