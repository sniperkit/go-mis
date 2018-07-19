package loanHistory

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(LoanHistory{}, []LoanHistory{})
}
