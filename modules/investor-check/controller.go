package investorCheck

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

type totalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

// FetchDatatables -  fetch data based on parameters sent by datatables
func FetchDatatables(ctx *iris.Context) {
	investors := []InvestorCheck{}
	totalData := totalData{}

	query := "SELECT id, name, \"idCardNo\", \"idCardFilename\", \"taxCardNo\", \"taxCardFilename\" "
	query += "FROM cif "
	query += "WHERE \"isValidated\" = false "
	query += "AND \"deletedAt\" IS NULL "

	queryTotalData := "SELECT count(cif.*) as \"totalRows\" "
	queryTotalData += "FROM cif "
	queryTotalData += "WHERE \"isValidated\" = false "
	queryTotalData += "AND \"deletedAt\" IS NULL "

	if ctx.URLParam("search") != "" {
		query += "AND name ~* '" + ctx.URLParam("search") + "' "
		queryTotalData += "AND name ~* '" + ctx.URLParam("search") + "' "
	}

	services.DBCPsql.Raw(query).Scan(&investors)
	services.DBCPsql.Raw(queryTotalData).Scan(&totalData)

	if ctx.URLParam("limit") != "" {
		query += "LIMIT " + ctx.URLParam("limit") + " "
	} else {
		query += "LIMIT 10 "
	}

	if ctx.URLParam("page") != "" {
		query += "OFFSET " + ctx.URLParam("page")
	} else {
		query += "OFFSET 0 "
	}

	services.DBCPsql.Raw(query).Scan(&investors)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      investors,
	})
}

// Verify - verify the selected investor
func Verify(ctx *iris.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	services.DBCPsql.Table("cif").Where("id = ?", id).Update("isValidated", true)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
	})
}
