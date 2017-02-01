package cif

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Cif{})
	services.BaseCrudInit(Cif{}, []Cif{})
}

func FetchAll(ctx *iris.Context) {
	cifInvestorBorrower := []CifInvestorBorrower{}

	query := "SELECT cif.\"id\", cif.\"cifNumber\", cif.\"name\", cif.\"isActivated\", cif.\"isValidated\", "
	query += "r_cif_borrower.\"borrowerId\" as \"borrowerId\", r_cif_investor.\"investorId\" as \"investorId\", "
	query += "r_cif_borrower.\"borrowerId\" IS NOT NULL as \"isBorrower\", r_cif_investor.\"investorId\" IS NOT NULL as \"isInvestor\" "
	query += "FROM cif "
	query += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"cifId\" = cif.\"id\" "
	query += "LEFT JOIN r_cif_investor ON r_cif_investor.\"cifId\" = cif.\"id\" "

	services.DBCPsql.Raw(query).Find(&cifInvestorBorrower)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   cifInvestorBorrower,
	})
}

func GetCifSummary(ctx *iris.Context) {
	cifSummary := CifSummary{}

	query := "SELECT (SELECT COUNT(cif.*) AS \"totalRegisteredCif\" FROM cif), (SELECT COUNT(investor.*) AS \"totalInvestor\" FROM investor), (SELECT COUNT(borrower.*) AS \"totalBorrower\" FROM borrower)"

	services.DBCPsql.Raw(query).Find(&cifSummary)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   cifSummary,
	})
}
