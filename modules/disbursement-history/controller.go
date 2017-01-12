package disbursementHistory

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&DisbursementHistory{})
	services.BaseCrudInit(DisbursementHistory{}, []DisbursementHistory{})
}
