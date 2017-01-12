package adjustment

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Adjustment{})
	services.BaseCrudInit(Adjustment{}, []Adjustment{})
}
