package investor

import (
	"fmt"

	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/virtual-account"
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
	query += "AND cif.\"deletedAt\" IS NULL AND cif.\"isValidated\" = true "
	query += "AND virtual_account.\"deletedAt\" IS NULL "

	investorWithoutVaSchema := []InvestorWithoutVaSchema{}
	services.DBCPsql.Raw(query).Scan(&investorWithoutVaSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   investorWithoutVaSchema,
	})
}

type InvestorVASchema struct {
	InvestorID   uint64 `json:"investorId"`
	InvestorNo   uint64 `json:"investorNo"`
	InvestorName string `json:"investorName"`
	VaBri        string `json:"vaBri"`
	VaBca        string `json:"vaBca"`
}

// InvestorRegisterVA - register VA to investor
func InvestorRegisterVA(ctx *iris.Context) {
	investorVASchema := InvestorVASchema{}

	if err := ctx.ReadJSON(&investorVASchema); err != nil {
		fmt.Println(investorVASchema)
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	vaSchemaBCA := virtualAccount.VirtualAccount{}
	vaSchemaBCA.VirtualAccountCode = "04435"
	vaSchemaBCA.VirtualAccountName = investorVASchema.InvestorName
	vaSchemaBCA.VirtualAccountNo = investorVASchema.VaBca

	services.DBCPsql.Create(vaSchemaBCA)

	rInvestorVaBca := &r.RInvestorVirtualAccount{InvestorId: investorVASchema.InvestorID, VirtualAccountId: vaSchemaBCA.ID}
	services.DBCPsql.Create(rInvestorVaBca)

	vaSchemaBRI := virtualAccount.VirtualAccount{}
	vaSchemaBRI.VirtualAccountCode = "99959"
	vaSchemaBRI.VirtualAccountName = investorVASchema.InvestorName
	vaSchemaBRI.VirtualAccountNo = investorVASchema.VaBri

	services.DBCPsql.Create(vaSchemaBRI)

	rInvestorVaBri := &r.RInvestorVirtualAccount{InvestorId: investorVASchema.InvestorID, VirtualAccountId: vaSchemaBRI.ID}
	services.DBCPsql.Create(rInvestorVaBri)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"bca": vaSchemaBCA,
			"bri": vaSchemaBRI,
		},
	})
}
