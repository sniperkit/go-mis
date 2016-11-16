package investor

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Investor{})
	services.BaseCrudInit(Investor{}, []Investor{})
}
