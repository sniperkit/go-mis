package agent

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&Agent{})
	services.BaseCrudInit(Agent{}, []Agent{})
}
