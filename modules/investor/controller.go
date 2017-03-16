package investor

import (
	"fmt"
	"strconv"

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

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

type InvestorSchema struct {
	AccountID  uint64 `gorm:"column:accountId" json:"accountId"`
	CifID      uint64 `gorm:"column:cifId" json:"cifId"`
	InvestorID uint64 `gorm:"column:investorId" json:"investorId"`
	Fullname   string `gorm:"column:fullname" json:"fullname"`
	Username   string `gorm:"column:username" json:"username"`
	PhoneNo    string `gorm:"column:phoneNo" json:"phoneNo"`
}

func GetInvestorForTopup(ctx *iris.Context) {
	queryCount := "SELECT count(*) as \"totalRows\" "
	queryCount += "FROM investor "
	queryCount += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.id "
	queryCount += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "
	queryCount += "JOIN r_account_investor ON r_account_investor.\"investorId\" = investor.id "
	queryCount += "WHERE cif.\"deletedAt\" IS NULL AND investor.\"deletedAt\" IS NULL "

	queryGetInvestor := "SELECT r_account_investor.id AS \"accountId\", cif.id AS \"cifId\", investor.id AS \"investorId\", cif.\"name\" AS \"fullname\", cif.\"username\" AS \"username\", cif.\"phoneNo\"   "
	queryGetInvestor += "FROM investor "
	queryGetInvestor += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.id "
	queryGetInvestor += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "
	queryGetInvestor += "JOIN r_account_investor ON r_account_investor.\"investorId\" = investor.id "
	queryGetInvestor += "WHERE cif.\"deletedAt\" IS NULL AND investor.\"deletedAt\" IS NULL "

	if ctx.URLParam("search") != "" {
		queryCount += "AND cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
		queryGetInvestor += "AND cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
	}

	var limitPagination int64 = 10
	var offset int64 = 0

	if ctx.URLParam("limit") != "" {
		queryGetInvestor += "LIMIT " + ctx.URLParam("limit") + " "
	} else {
		queryGetInvestor += "LIMIT " + strconv.FormatInt(limitPagination, 10) + " "
	}

	if ctx.URLParam("page") != "" {
		offset, _ = strconv.ParseInt(ctx.URLParam("page"), 10, 64)
		queryGetInvestor += "OFFSET " + strconv.FormatInt(offset, 10)
	} else {
		queryGetInvestor += "OFFSET 0"
	}

	totalData := TotalData{}
	services.DBCPsql.Raw(queryCount).Find(&totalData)

	investorSchema := []InvestorSchema{}
	services.DBCPsql.Raw(queryGetInvestor).Find(&investorSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      investorSchema,
	})
}
