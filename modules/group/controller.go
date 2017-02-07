package group

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Group{})
	services.BaseCrudInit(Group{}, []Group{})
}

// FetchAll - fetchAll group data
func FetchAll(ctx *iris.Context) {
	groupBranchAreaAgent := []GroupBranchAreaAgent{}

	query := "SELECT \"group\".\"id\", \"group\".\"name\", \"group\".\"createdAt\", branch.\"name\" as \"branch\", area.\"name\" as \"area\", agent.\"fullname\" as \"agent\" "
	query += "FROM \"group\" "
	query += "LEFT JOIN r_group_agent ON r_group_agent.\"groupId\" = \"group\".\"id\" "
	query += "LEFT JOIN agent ON agent.\"id\" = r_group_agent.\"agentId\" "
	query += "LEFT JOIN r_group_branch ON r_group_branch.\"groupId\" = \"group\".\"id\" "
	query += "LEFT JOIN branch ON branch.\"id\" = \"r_group_branch\".\"branchId\" "
	query += "LEFT JOIN r_area_branch ON r_group_branch.\"branchId\" = r_area_branch.\"branchId\" "
	query += "LEFT JOIN area ON r_area_branch.\"areaId\" = area.\"id\" "
	query += "WHERE \"group\".\"deletedAt\" is NULL"

	if e := services.DBCPsql.Raw(query).Find(&groupBranchAreaAgent).Error; e != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "failed",
			"data":   e,
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   groupBranchAreaAgent,
	})
}
