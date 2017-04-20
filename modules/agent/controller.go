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
	var tmp uint64 // why this? 0 is untype constant
	tmp = 0
	agentSchema := []Agent{}

	// query := "SELECT agent.* "
	// query += "FROM agent "
	// query += "INNER JOIN r_branch_agent ON r_branch_agent.\"agentId\" = agent.id "
	// query += "WHERE r_branch_agent.\"branchId\" = ? AND agent.\"deletedAt\" IS NULL"

	query := ""
	// if not superadmin
	// TODO: use role instead of branchID
	if branchID != tmp {
		query += "SELECT agent.id, agent.\"picUrl\", agent.\"username\", agent.fullname, agent.address,  "
		query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = '0' AND kecamatan = '0' AND kelurahan = '0' LIMIT 1) AS \"province\", "
		query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = agent.city AND kecamatan = '0' AND kelurahan = '0' LIMIT 1) AS \"city\", "
		query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = agent.city AND kecamatan = agent.kecamatan AND kelurahan = '0' LIMIT 1) AS \"kecamatan\", "
		query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = agent.city AND kecamatan = agent.kecamatan AND kelurahan = agent.kelurahan LIMIT 1) AS \"kelurahan\" "
		query += "FROM agent "
		query += "INNER JOIN r_branch_agent ON r_branch_agent.\"agentId\" = agent.id "
		query += "WHERE r_branch_agent.\"branchId\" = ? AND agent.\"deletedAt\" IS NULL"

		services.DBCPsql.Raw(query).Scan(&agentSchema)
	} else {
			query += `SELECT agent.id, agent."picUrl", agent."username", agent.fullname, agent.address,
			(SELECT "name" FROM inf_location 
			WHERE province = agent.province 
			AND city = '0' 
			AND kecamatan = '0' 
			AND kelurahan = '0' LIMIT 1) AS "province",
			(SELECT "name" FROM inf_location 
			WHERE province = agent.province 
			AND city = agent.city 
			AND kecamatan = '0' 
			AND kelurahan = '0' LIMIT 1) AS "city",
			(SELECT "name" FROM inf_location 
			WHERE province = agent.province 
			AND city = agent.city 
			AND kecamatan = agent.kecamatan 
			AND kelurahan = '0' LIMIT 1) AS "kecamatan",
			(SELECT "name" FROM inf_location 
			WHERE province = agent.province 
			AND city = agent.city 
			AND kecamatan = agent.kecamatan 
			AND kelurahan = agent.kelurahan LIMIT 1) AS "kelurahan"
			from agent
			INNER JOIN r_branch_agent ON r_branch_agent."agentId" = agent.id
			WHERE agent."deletedAt" IS NULL` 

		services.DBCPsql.Raw(query).Scan(&agentSchema)
	}


	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   agentSchema,
	})
}
