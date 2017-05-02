package prospectiveBorrower

import (
	"time"

	"bitbucket.org/go-mis/modules/branch"
	"bitbucket.org/go-mis/modules/role"
	"bitbucket.org/go-mis/modules/survey"
	"bitbucket.org/go-mis/modules/user-mis"
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
	roles := []role.Role{}
	qRole := `SELECT "role".* FROM "role"
		JOIN r_user_mis_role rumr ON rumr."roleId" = "role".id 
		WHERE ("role"."name" ~* 'admin' OR "role"."name" ~* 'area') AND rumr."userMisId" = ?
	`

	userMis := ctx.Get("USER_MIS").(userMis.UserMis)
	services.DBCPsql.Raw(qRole, userMis.ID).Scan(&roles)

	surveySchema := []SurveySchema{}
	if len(roles) == 0 {
		branchID := ctx.Get("BRANCH_ID")
		q := `SELECT survey.*, branch."name" as "branch", "group"."name" as "group", agent.fullname as "agent"
			FROM survey 
			LEFT JOIN r_group_branch ON r_group_branch."groupId" = survey."groupId"
			LEFT JOIN branch ON branch.id = r_group_branch."branchId"
			LEFT JOIN agent ON agent.id = survey."agentId"
			LEFT JOIN "group" ON "group".id = survey."groupId"
			WHERE branch.id = ? AND "isMigrate" = false AND "isApprove" = false AND survey."deletedAt" IS NULL
			ORDER BY survey.id asc
		`
		services.DBCPsql.Table("survey").Raw(q, branchID).Scan(&surveySchema)
	} else {
		qBranch := `
		SELECT branch.* FROM branch
		JOIN r_area_branch rab ON rab."branchId" = branch.id
		JOIN r_area_user_mis raum ON raum."areaId" = rab."areaId"
		WHERE raum."userMisId" = ?
		`
		branches := []branch.Branch{}
		services.DBCPsql.Raw(qBranch, userMis.ID).Scan(&branches)

		branchIds := make([]uint64, len(branches))
		for i, current := range branches {
			branchIds[i] = current.ID
		}

		q := `SELECT survey.*, branch."name" as "branch", "group"."name" as "group", agent.fullname as "agent"
			FROM survey 
			LEFT JOIN r_group_branch ON r_group_branch."groupId" = survey."groupId"
			LEFT JOIN branch ON branch.id = r_group_branch."branchId"
			LEFT JOIN agent ON agent.id = survey."agentId"
			LEFT JOIN "group" ON "group".id = survey."groupId"
			WHERE branch.id in (?) AND "isMigrate" = false AND "isApprove" = false AND survey."deletedAt" IS NULL
			ORDER BY survey.id asc
		`
		services.DBCPsql.Table("survey").Raw(q, branchIds).Scan(&surveySchema)
	}

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
