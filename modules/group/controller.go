package group

import (
	"bitbucket.org/go-mis/modules/branch"
	"bitbucket.org/go-mis/modules/role"
	"bitbucket.org/go-mis/modules/user-mis"
	"bitbucket.org/go-mis/services"
	"bitbucket.org/go-mis/modules/r"
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

func GroupDetail(ctx *iris.Context){

	id := ctx.Get("id");
	groupBorrower := []GroupAgentBorrower{}

	query := "SELECT \"group\".\"id\", \"group\".\"name\" as \"name\", \"group\".\"lat\" as \"lat\",\"group\".\"lng\" as \"lng\", \"group\".\"scheduleDay\" as \"scheduleDay\", \"group\".\"scheduleTime\" as \"scheduleTime\", \"group\".\"name\", cif.\"name\" as \"borrowerName\" "
	query += "FROM \"group\" "
	query += "LEFT JOIN r_group_borrower rgb ON rgb.\"groupId\" = \"group\".\"id\" "
	query += "LEFT JOIN borrower ON borrower.\"id\" = rgb.\"borrowerId\" "
	query += "LEFT JOIN r_cif_borrower rcb ON rcb.\"borrowerId\" = borrower.\"id\" "
	query += "LEFT JOIN cif ON cif.\"id\" = rcb.\"cifId\" WHERE \"group\".\"id\" = ? "

	if e := services.DBCPsql.Raw(query, id).Scan(&groupBorrower).Error; e != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "failed",
			"data":   e,
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   groupBorrower,
	})
}

func Create(ctx *iris.Context){

	type Payload struct {
		ID 				uint64 		`json:"_id"`
		Name 			string 		`json:"name"`
		Lat 			float64 	`json:"lat"`
		Lng 			float64 	`json:"lng"`
		Agent 		uint64 		`json:"agentId"`
		Branch 		uint64 		`json:"branchId"`
	}

	m := Payload{}
	err := ctx.ReadJSON(&m)

	g := Group{}
	g.Name = m.Name;
	g.Lat = m.Lat;
	g.Lng = m.Lng;

	if err != nil { 
		panic(err) 
	}else{
		services.DBCPsql.Create(&g);

		rga := r.RGroupAgent{}
		rga.GroupId = g.ID;
		rga.AgentId = m.Agent;

		if err := services.DBCPsql.Create(&rga).Error; err != nil {
			panic(err)
		}

		rgb := r.RGroupBranch{}
		rgb.GroupId = g.ID;
		rgb.BranchId = m.Branch;

		services.DBCPsql.Create(&rgb);

	}
	ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": m})

}
