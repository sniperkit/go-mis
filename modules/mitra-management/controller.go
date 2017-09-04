package mitramanagement

import (
	"log"

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

func GetStatusAll(ctx *iris.Context) {
	s := []Status{}
	err := services.DBCPsql.Where("type = 'mitra_management'").Find(&s).Error
	if err != nil {
		log.Println("[INFO] Params is not valid")
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": err.Error(),
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   s,
	})
}
