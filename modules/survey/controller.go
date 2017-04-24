package survey

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Survey{})
	services.BaseCrudInit(Survey{}, []Survey{})
}

// GetProspectiveBorrower - get prospective borrower which is not migrated
func GetProspectiveBorrower(ctx *iris.Context) {
	aFields := []AFields{}
	query := "SELECT * FROM a_fields WHERE \"is_migrated\" = false OR is_migrated IS NULL"
	services.DBCPsqlSurvey.Raw(query).Find(&aFields)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   aFields,
	})
}

// GetProspectiveBorrowerArchived - get rejected prospective borrower
func GetProspectiveBorrowerArchived(ctx *iris.Context) {
	aFields := []AFields{}
	query := "SELECT * FROM a_fields WHERE is_migrated = true AND is_approve = false"
	services.DBCPsqlSurvey.Raw(query).Find(&aFields)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   aFields,
	})
}

// GetProspectiveBorrowerDetail - get prospective borrower detail
func GetProspectiveBorrowerDetail(ctx *iris.Context) {
	aFields := []AFields{}
	// query := "SELECT * FROM a_fields WHERE (\"is_migrated\" = false OR is_migrated IS NULL) AND \"answer_id\" = ?"
	query := "SELECT * FROM a_fields WHERE \"answer_id\" = ?"
	query += " ORDER BY \"key\" ASC"
	services.DBCPsqlSurvey.Raw(query, ctx.Param("id")).Find(&aFields)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   aFields,
	})
}
