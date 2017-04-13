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
  where cif.name ~* ? 
  and investor."isInstitutional" = true
  and not exists (
  	select 1 from r_investor_product_pricing ripp
  	where ripp."investorId" = investor.id
  )`

	services.DBCPsql.Raw(query, searchStr).Scan(&sInv)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": sInv,
	})
}

func GetInvestorsByProuctPricing (ctx *iris.Context) {
	ppId := ctx.Param("id")
	sInv := []InvestorSearch{}

	query := `select investor.id, cif.name
	from product_pricing
	join r_investor_product_pricing ripp on ripp."productPricingId" = product_pricing.id
	join investor on investor.id = ripp."investorId"
	join r_cif_investor rci on rci."investorId" = investor.id
	join cif on rci."cifId" = cif.id
	where ripp."productPricingId" = ?`

	services.DBCPsql.Raw(query, ppId).Scan(&sInv)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": sInv,
	})
}
