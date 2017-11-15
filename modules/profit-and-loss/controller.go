package profitAndLoss

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&ProfitAndLoss{})
	services.BaseCrudInit(ProfitAndLoss{}, []ProfitAndLoss{})
}
