package productPricing

import (
	"bitbucket.org/go-mis/services"
	"bitbucket.org/go-mis/modules/r"
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
  where cif.name LIKE '%?%'
  and investor."isInstitutional" = true
  and not exists (
  	select 1 from r_investor_product_pricing ripp
  	where ripp."investorId" = investor.id and ripp."deletedAt" is null
  )`

	services.DBCPsql.Raw(query, searchStr).Scan(&sInv)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   sInv,
	})
}

func Create(ctx *iris.Context) {
	m := ProductPricing{}

	

	err := ctx.ReadJSON(&m)
	if err != nil {
		panic(err)
	}
	if *m.IsInstutitional == false {
		pplist := []ProductPricing{}
		query := `select * from product_pricing where product_pricing."isInstitutional" = false AND product_pricing."deletedAt" IS NULL `
		query += `and ( product_pricing."startDate" between ? and ? or product_pricing."endDate" between ? and ? or `
		query += `? between product_pricing."startDate" and product_pricing."endDate" or ? between product_pricing."startDate" and product_pricing."endDate") `

		services.DBCPsql.Raw(query, m.StartDate, m.EndDate, m.StartDate, m.EndDate, m.StartDate, m.EndDate).Scan(&pplist)

		if len(pplist) > 0 {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{"status": "error", "message": "Date Overlap, please choose another date.", "data": pplist})
		} else {
			services.DBCPsql.Create(&m)
			ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": m})
		}
	} else {
		services.DBCPsql.Create(&m)

		for _, val := range m.Investors {
			r := r.RInvestorProductPricing{}
			r.InvestorId=val.ID
			r.ProductPricingId=m.ID
			if err := services.DBCPsql.Create(&r).Error; err != nil {
				panic(err)
			}
		}

		ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": m})

	}
}

func GetInvestorsByProductPricing(ctx *iris.Context) {
	ppId := ctx.Param("id")
	sInv := []InvestorSearchByProductPricing{}

	query := `select investor.id, cif.name, ripp.id as "rippId"
	from product_pricing
	join r_investor_product_pricing ripp on ripp."productPricingId" = product_pricing.id
	join investor on investor.id = ripp."investorId"
	join r_cif_investor rci on rci."investorId" = investor.id
	join cif on rci."cifId" = cif.id
	where ripp."productPricingId" = ? and ripp."deletedAt" is null`

	services.DBCPsql.Raw(query, ppId).Scan(&sInv)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   sInv,
	})
}

func DeleteProductPricing(ctx *iris.Context){
	ppId := ctx.Param("id")
	productPricing := ProductPricing{}

	query := "WITH upd AS (UPDATE product_pricing SET \"deletedAt\" = now() WHERE id = ? RETURNING id) "
	query += "UPDATE r_investor_product_pricing ripp SET \"deletedAt\" = now() FROM upd WHERE ripp.\"productPricingId\" = upd.\"id\""

	services.DBCPsql.Raw(query, ppId).Scan(&productPricing);

	ctx.JSON(iris.StatusOK, iris.Map{"data": productPricing})

}
func DeleteProductPricingInvestor(ctx *iris.Context){
	ppId := ctx.Param("ppId")
	invId := ctx.Param("invId")
	productPricing := ProductPricing{}
	query := `UPDATE r_investor_product_pricing SET "deletedAt" = now() WHERE "investorId" = ? AND "productPricingId" = ?`
	services.DBCPsql.Raw(query, invId, ppId).Scan(&productPricing);
	ctx.JSON(iris.StatusOK, iris.Map{"data": productPricing})
}
