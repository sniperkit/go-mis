package agent

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Agent{})
	services.BaseCrudInit(Agent{}, []Agent{})
}

func GetAllAgentByBranchID(ctx *iris.Context) {
	branchID := ctx.Get("BRANCH_ID")
	agentSchema := []Agent{}

	query := "SELECT agent.* "
	query += "FROM agent "
	query += "INNER JOIN r_branch_agent ON r_branch_agent.\"agentId\" = agent.id "
	query += "WHERE r_branch_agent.\"branchId\" = ? "

	services.DBCPsql.Raw(query, branchID).Scan(&agentSchema)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   agentSchema,
	})
}
