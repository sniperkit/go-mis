package sector

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Sector{})
	services.BaseCrudInit(Sector{}, []Sector{})
}
