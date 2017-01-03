package area

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Area{})
	services.BaseCrudInit(Area{}, []Area{})
}
