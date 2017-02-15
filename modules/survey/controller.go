package survey

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GetProspectiveBorrower - get prospective borrower which si not migrated
func GetProspectiveBorrower(ctx *iris.Context) {
	aFields := []AFields{}
	query := "SELECT * FROM a_fields WHERE \"is_migrated\" = false"
	services.DBCPsqlSurvey.Raw(query).Find(&aFields)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   aFields,
	})
}

// GetProspectiveBorrowerDetail - get prospective borrower detail
func GetProspectiveBorrowerDetail(ctx *iris.Context) {
	aFields := []AFields{}
	query := "SELECT * FROM a_fields WHERE \"is_migrated\" = false AND \"answer_id\" = ?"
	query += " ORDER BY \"key\" ASC"
	services.DBCPsqlSurvey.Raw(query, ctx.Param("id")).Find(&aFields)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   aFields,
	})
}
