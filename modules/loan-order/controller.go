package loanOrder

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&LoanOrder{})
	services.BaseCrudInit(LoanOrder{}, []LoanOrder{})
}
