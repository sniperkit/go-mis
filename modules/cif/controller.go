package cif

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Cif{})
	services.BaseCrudInit(Cif{}, []Cif{})
}

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

func FetchAll(ctx *iris.Context) {
	totalData := TotalData{}

	queryTotalData := "SELECT count(cif.*) as \"totalRows\" "
	queryTotalData += "FROM cif "
	queryTotalData += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"cifId\" = cif.\"id\" "
	queryTotalData += "LEFT JOIN r_cif_investor ON r_cif_investor.\"cifId\" = cif.\"id\" "

	if ctx.URLParam("search") != "" {
		// queryTotalData += "WHERE cif.name LIKE '%" + ctx.URLParam("search") + "%' "
		queryTotalData += "WHERE cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
	}

	services.DBCPsql.Raw(queryTotalData).Find(&totalData)

	var limitPagination int64 = 10
	var offset int64 = 0
	cifInvestorBorrower := []CifInvestorBorrower{}

	query := "SELECT cif.\"id\", cif.\"cifNumber\", cif.\"name\", cif.\"isActivated\", cif.\"isValidated\", "
	query += "r_cif_borrower.\"borrowerId\" as \"borrowerId\", r_cif_investor.\"investorId\" as \"investorId\", "
	query += "r_cif_borrower.\"borrowerId\" IS NOT NULL as \"isBorrower\", r_cif_investor.\"investorId\" IS NOT NULL as \"isInvestor\" "
	query += "FROM cif "
	query += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"cifId\" = cif.\"id\" "
	query += "LEFT JOIN r_cif_investor ON r_cif_investor.\"cifId\" = cif.\"id\" "

	if ctx.URLParam("search") != "" {
		// query += "WHERE cif.name LIKE '%" + ctx.URLParam("search") + "%' "
		query += "WHERE cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
	}

	if ctx.URLParam("limit") != "" {
		query += "LIMIT " + ctx.URLParam("limit") + " "
	} else {
		query += "LIMIT " + strconv.FormatInt(limitPagination, 10) + " "
	}

	if ctx.URLParam("page") != "" {
		offset, _ = strconv.ParseInt(ctx.URLParam("page"), 10, 64)
		query += "OFFSET " + strconv.FormatInt(offset, 10)
	} else {
		query += "OFFSET 0"
	}

	services.DBCPsql.Raw(query).Find(&cifInvestorBorrower)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      cifInvestorBorrower,
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

type CifInvestorAccount struct {
	CifID        uint64  `gorm:"column:cifId" json:"cifId"`
	Name         string  `gorm:"column:name" json:"name"`
	TotalBalance float64 `gorm:"column:totalBalance" json:"totalBalance"`
}

// GetCifInvestorAccount - Get CIF, investor and account data
func GetCifInvestorAccount(ctx *iris.Context) {
	email := ctx.URLParam("email")

	query := "SELECT cif.id AS \"cifId\", cif.name, account.\"totalBalance\" "
	query += "FROM cif "
	query += "JOIN r_cif_investor ON r_cif_investor.\"cifId\" = cif.id "
	query += "JOIN r_account_investor ON r_account_investor.\"investorId\" = r_cif_investor.\"investorId\" "
	query += "JOIN account ON account.id = r_account_investor.\"accountId\" "
	query += "WHERE cif.username = ? "

	cifInvestorAccountObj := new(CifInvestorAccount)

	services.DBCPsql.Raw(query, email).Scan(&cifInvestorAccountObj)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   cifInvestorAccountObj,
	})
}
