package mitramanagement

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Status{})
	services.DBCPsql.AutoMigrate(&Reason{})
}

// GetPortfolioAtRisk - Get loan data where stage is equal to 'PAR'
func GetPortfolioAtRisk(ctx *iris.Context) {

}

func FindPortfolioAtRisk(parList []PortfolioAtRisk) (uint64, error) {
	return 0, nil
}
