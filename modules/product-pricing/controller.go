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

	query := `select investor.id, cif.name
	from investor
	join r_cif_investor on r_cif_investor."investorId" = investor.id
	join cif on r_cif_investor."cifId" = cif.id
	where cif.name like ?`

	services.DBCPsql.Raw(query, "%"+searchStr+"%").Scan(&sInv)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": sInv,
	})
}
