package loanHistory

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&LoanHistory{})
	services.BaseCrudInit(LoanHistory{}, []LoanHistory{})
}
