package installmentHistory

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&InstallmentHistory{})
	services.BaseCrudInit(InstallmentHistory{}, []InstallmentHistory{})
}
