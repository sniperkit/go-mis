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
	groups := []GroupAreaAgent{}

	query := "SELECT \"group\".\"id\", \"group\".\"name\", \"group\".\"createdAt\", area.\"name\" as \"areaName\", agent.\"fullname\" as \"agentName\" "
	query += "FROM \"group\" "
	query += "JOIN r_group_agent ON r_group_agent.\"groupId\" = \"group\".\"id\" "
	query += "JOIN agent ON agent.\"id\" = r_group_agent.\"agentId\" "
	query += "JOIN r_group_branch ON r_group_branch.\"groupId\" = \"group\".\"id\" "
	query += "JOIN r_area_branch ON r_group_branch.\"branchId\" = r_area_branch.\"branchId\" "
	query += "JOIN area ON r_area_branch.\"areaId\" = area.\"id\" "

	services.DBCPsql.Raw(query).Find(&groups)
	ctx.JSON(iris.StatusOK, iris.Map{"data": groups})
}
