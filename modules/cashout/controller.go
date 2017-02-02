package cashout

import (
	"fmt"
	"strings"

	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Cashout{})
	services.BaseCrudInit(Cashout{}, []Cashout{})
}

// FetchAll - fetch cashout data
func FetchAll(ctx *iris.Context) {
	cashoutInvestors := []CashoutInvestor{}
	stage := ctx.URLParam("stage")
	fmt.Print(stage)

	query := "SELECT cashout.*, "
	query += "cif.\"name\" as \"investorName\", "
	query += "account.\"totalBalance\" as \"totalBalance\" "
	// query += "r_cif_borrower.\"borrowerId\" IS NOT NULL as \"isBorrower\", r_cif_investor.\"investorId\" IS NOT NULL as \"isInvestor\" "
	query += "FROM cashout "
	query += "JOIN r_investor_cashout ON r_investor_cashout.\"cashoutId\" = cashout.\"id\" "
	query += "JOIN investor ON investor.\"id\" = r_investor_cashout.\"investorId\" "
	query += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.\"id\" "
	query += "JOIN cif ON r_cif_investor.\"cifId\" = cif.\"id\" "
	query += "JOIN r_account_investor ON r_account_investor.\"investorId\" = investor.\"id\" "
	query += "JOIN account ON r_account_investor.\"accountId\" = account.\"id\" "
	if len(strings.TrimSpace(stage)) == 0 {
		services.DBCPsql.Raw(query).Find(&cashoutInvestors)
	} else {
		query += "WHERE cashout.\"stage\" = ? "
		services.DBCPsql.Raw(query, stage).Find(&cashoutInvestors)
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   cashoutInvestors,
	})
}
