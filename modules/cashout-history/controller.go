package cashoutHistory

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&CashoutHistory{})
	services.BaseCrudInit(CashoutHistory{}, []CashoutHistory{})
}
