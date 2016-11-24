package borrowerProspective

import "bitbucket.org/go-mis/services"

func Init() {
	services.DBCPsql.AutoMigrate(&BorrowerProspective{})
	services.BaseCrudInit(BorrowerProspective{}, []BorrowerProspective{})
}
