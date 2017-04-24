package prospectiveBorrower

import (
	"time"

	"bitbucket.org/go-mis/modules/survey"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

type SurveySchema struct {
	SurveyID         uint64    `gorm:"column:id" json:"id"`
	AgentID          uint64    `gorm:"column:agentId" json:"agentId"`
	BranchID         uint64    `gorm:"column:branchId" json:"branchId"`
	GroupID          uint64    `gorm:"column:groupId" json:"groupId"`
	Branch           string    `gorm:"column:branch" json:"branch"`
	Group            string    `gorm:"column:group" json:"group"`
	Agent            string    `gorm:"column:agent" json:"agent"`
	Fullname         string    `gorm:"column:fullname" json:"fullname"`
	CreditScoreGrade string    `gorm:"column:creditScoreGrade" json:"creditScoreGrade"`
	CreditScoreValue float64   `gorm:"column:creditScoreValue" json:"creditScoreValue"`
	CreatedAt        time.Time `gorm:"column:createdAt" json:"createdAt"`
}

// GetProspectiveBorrower - Get prospective borrower v2
func GetProspectiveBorrower(ctx *iris.Context) {
	surveySchema := []SurveySchema{}
	q := `SELECT survey.*, branch."name" as "branch", "group"."name" as "group", agent.fullname as "agent"
	 	FROM survey 
		LEFT JOIN branch ON branch.id = survey."branchId"
		LEFT JOIN agent ON agent.id = survey."agentId"
		LEFT JOIN "group" ON "group".id = survey."groupId"
		WHERE "isMigrate" = false AND "isApprove" = false AND survey."deletedAt" IS NULL
		ORDER BY survey.id asc
	`
	services.DBCPsql.Table("survey").Raw(q).Scan(&surveySchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   surveySchema,
	})
}

// GetProspectiveBorrowerDetail - Get prospective borrower detail
func GetProspectiveBorrowerDetail(ctx *iris.Context) {
	surveySchema := survey.Survey{}

	services.DBCPsql.Table("survey").Where("id = ? AND \"deletedAt\" IS NULL ", ctx.Param("id")).Scan(&surveySchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   surveySchema,
	})
}

// UpdateStatusProspectiveBorrower - Update status
func UpdateStatusProspectiveBorrower(ctx *iris.Context) {
	id := ctx.Param("id")
	status := ctx.Param("status")

	services.DBCPsql.Table("survey").Where("id = ? AND \"deletedAt\" IS NULL", id).UpdateColumn("isMigrate", true)
	services.DBCPsql.Table("survey").Where("id = ? AND \"deletedAt\" IS NULL", id).UpdateColumn("isApprove", status)
	services.DBCPsql.Table("survey").Where("id = ? AND \"deletedAt\" IS NULL", id).UpdateColumn("updatedAt", "current_timestamp")

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}
