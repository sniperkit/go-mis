package installmentPresence

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&InstallmentPresence{})
	services.BaseCrudInit(InstallmentPresence{}, []InstallmentPresence{})
}
