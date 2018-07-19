package loanMonitoring

import "bitbucket.org/go-mis/services"

func Init() {
		services.BaseCrudInit(LoanMonitoring{}, []LoanMonitoring{})
}
