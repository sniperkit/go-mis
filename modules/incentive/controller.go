package incentive

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Incentive{})
	services.BaseCrudInit(Incentive{}, []Incentive{})
}
