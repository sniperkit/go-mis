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
	where cif.name ~* ? and investor."isInstitutional" = true`

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
	if *m.IsInstutitional == true {
		pplist := []ProductPricing{}
		query := `select * from product_pricing where product_pricing."isInstitutional" = true `
		query += `and ( product_pricing."startDate" between ? and ? or product_pricing."endDate" between ? and ? ) `

		services.DBCPsql.Raw(query, m.StartDate, m.EndDate, m.StartDate, m.EndDate).Scan(&pplist)

		if len(pplist) > 0 {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": "date overlap", "data": pplist})
		} else {
			services.DBCPsql.Table("product_pricing").Create(m)
			ctx.JSON(iris.StatusOK, iris.Map{"data": m})
		}

	} else {
		services.DBCPsql.Table("product_pricing").Create(m)

		ctx.JSON(iris.StatusOK, iris.Map{"data": m})

	}

}
