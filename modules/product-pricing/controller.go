package productPricing

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&ProductPricing{})
	services.BaseCrudInit(ProductPricing{}, []ProductPricing{})
}

func SearchInvestor(ctx *iris.Context) {
	searchStr := ctx.Param("searchStr")
	sInv := []InvestorSearch{}

	query := `select id, name from cif where name like ?`
	services.DBCPsql.Raw(query, "%"+searchStr+"%").Scan(&sInv)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": sInv,
	})
}
