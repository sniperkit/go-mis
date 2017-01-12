package loanMonitoring

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&LoanMonitoring{})
	services.BaseCrudInit(LoanMonitoring{}, []LoanMonitoring{})
}
