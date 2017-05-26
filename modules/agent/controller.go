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

	// query := "SELECT agent.* "
	// query += "FROM agent "
	// query += "INNER JOIN r_branch_agent ON r_branch_agent.\"agentId\" = agent.id "
	// query += "WHERE r_branch_agent.\"branchId\" = ? AND agent.\"deletedAt\" IS NULL"

	query := ""
	// if not superadmin
	// TODO: use role instead of branchID
	if branchID.(uint64) > 0 {
		query += "SELECT agent.id, agent.\"picUrl\", agent.\"username\", agent.fullname, agent.address,  "
		query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = '0' AND kecamatan = '0' AND kelurahan = '0' LIMIT 1) AS \"province\", "
		query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = agent.city AND kecamatan = '0' AND kelurahan = '0' LIMIT 1) AS \"city\", "
		query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = agent.city AND kecamatan = agent.kecamatan AND kelurahan = '0' LIMIT 1) AS \"kecamatan\", "
		query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = agent.city AND kecamatan = agent.kecamatan AND kelurahan = agent.kelurahan LIMIT 1) AS \"kelurahan\" "
		query += "FROM agent "
		query += "INNER JOIN r_branch_agent ON r_branch_agent.\"agentId\" = agent.id "
		query += "WHERE r_branch_agent.\"branchId\" = ? AND agent.\"deletedAt\" IS NULL"

		services.DBCPsql.Raw(query, branchID).Scan(&agentSchema)
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

func GetAgentById(ctx *iris.Context){
		result := AgentBranch{}
		query := "SELECT agent.\"id\" AS id, "
	  query += "agent.\"username\" AS \"username\", "
	  query += "agent.\"fullname\" AS \"fullname\", "
	  query += "agent.\"password\" AS \"password\", "
	  query += "agent.\"bankName\" AS \"bankName\", "
		query += "agent.\"bankAccountName\" AS \"bankAccountName\", "
	  query += "agent.\"bankAccountNo\" AS \"bankAccountName\", "
	  query += "agent.\"picUrl\" AS \"picUrl\", "
	  query += "agent.\"phoneNo\" AS \"phoneNo\", "
	  query += "agent.\"address\" AS \"address\", "
	  query += "agent.\"kelurahan\" AS \"kelurahan\", "
	  query += "agent.\"kecamatan\" AS \"kecamatan\", "
	  query += "agent.\"city\" AS \"city\", "
		query += "agent.\"province\" AS \"province\", "
		query += "agent.\"nationality\" AS \"nationality\", "
	  query += "agent.\"lat\" AS \"lat\", "
	  query += "agent.\"lng\" AS \"lng\", "
	  query += "branch.\"name\" AS \"branchName\" "
		query += "FROM agent "
		query += "LEFT JOIN r_branch_agent ON r_branch_agent.\"agentId\" = agent.\"id\" "
		query += "LEFT JOIN branch ON branch.\"id\" = r_branch_agent.\"branchId\" "
		query += "WHERE agent.\"id\" = ?"

		id := ctx.Get("id")
		services.DBCPsql.Raw(query, id).Scan(&result)
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   result,
		})

}
