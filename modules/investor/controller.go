package investor

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Investor{})
	services.BaseCrudInit(Investor{}, []Investor{})
}

// InvestorWithoutVA - Retrieve list of investor without VA
func InvestorWithoutVA(ctx *iris.Context) {
	query := "SELECT cif.\"name\", investor.id AS \"investorId\", investor.\"investorNo\", virtual_account.\"bankName\", virtual_account.\"virtualAccountNo\", virtual_account.\"virtualAccountName\" "
	query += "FROM investor "
	query += "LEFT OUTER JOIN r_investor_virtual_account ON r_investor_virtual_account.\"investorId\" = investor.id  "
	query += "LEFT OUTER JOIN virtual_account ON virtual_account.id = r_investor_virtual_account.\"investorId\" "
	query += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.id "
	query += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "
	query += "WHERE virtual_account.\"virtualAccountNo\" IS NULL  "
	query += "AND virtual_account.\"deletedAt\" IS NULL "

	investorWithoutVaSchema := []InvestorWithoutVaSchema{}
	services.DBCPsql.Raw(query).Scan(&investorWithoutVaSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   investorWithoutVaSchema,
	})
}
