package group

import (
	"bitbucket.org/go-mis/modules/branch"
	"bitbucket.org/go-mis/modules/role"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Group{})
	services.BaseCrudInit(Group{}, []Group{})
}

// FetchAll - fetchAll group data
func FetchAll(ctx *iris.Context) {
	roles := []role.Role{}
	qRole := `SELECT "role".* FROM "role"
		JOIN r_user_mis_role rumr ON rumr."roleId" = "role".id 
		WHERE ("role"."name" ~* 'admin' OR "role"."name" ~* 'area') AND rumr."userMisId" = ?
	`

	userMis := ctx.Get("USER_MIS").(userMis.UserMis)
	services.DBCPsql.Raw(qRole, userMis.ID).Scan(&roles)

	groupBranchAreaAgent := []GroupBranchAreaAgent{}
	query := "SELECT \"group\".\"id\", \"group\".\"name\", \"group\".\"createdAt\", branch.\"name\" as \"branch\", area.\"name\" as \"area\", agent.\"fullname\" as \"agent\" "
	query += "FROM \"group\" "
	query += "LEFT JOIN r_group_agent ON r_group_agent.\"groupId\" = \"group\".\"id\" "
	query += "LEFT JOIN agent ON agent.\"id\" = r_group_agent.\"agentId\" "
	query += "LEFT JOIN r_group_branch ON r_group_branch.\"groupId\" = \"group\".\"id\" "
	query += "LEFT JOIN branch ON branch.\"id\" = \"r_group_branch\".\"branchId\" "
	query += "LEFT JOIN r_area_branch ON r_group_branch.\"branchId\" = r_area_branch.\"branchId\" "
	query += "LEFT JOIN area ON r_area_branch.\"areaId\" = area.\"id\" "
	if len(roles) == 0 {
		query += "WHERE \"group\".\"deletedAt\" is NULL AND branch.id = ?"

		branchID := ctx.Get("BRANCH_ID")
		if e := services.DBCPsql.Raw(query, branchID).Scan(&groupBranchAreaAgent).Error; e != nil {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status": "failed",
				"data":   e,
			})
			return
		}
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

		query += "WHERE \"group\".\"deletedAt\" is NULL AND branch.id in (?) "
		if e := services.DBCPsql.Raw(query, branchIds).Scan(&groupBranchAreaAgent).Error; e != nil {
			ctx.JSON(iris.StatusInternalServerError, iris.Map{
				"status": "failed",
				"data":   e,
			})
			return
		}
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   groupBranchAreaAgent,
	})
}
