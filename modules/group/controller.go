package group

import (
	"bitbucket.org/go-mis/modules/branch"
	"bitbucket.org/go-mis/modules/r"
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

func SearchGroup(ctx *iris.Context){

	searchStr := ctx.Param("searchStr")
	branchId := ctx.Param("branchId")
	sGroup := []GroupSearch{}

	query := `SELECT "group".id,"group"."name" FROM "group" 
JOIN r_loan_group ON r_loan_group."groupId"="group".id
JOIN r_loan_branch ON r_loan_branch."loanId"=r_loan_group."loanId"
WHERE r_loan_branch."branchId"=? AND "group"."name" ILIKE ?
GROUP BY "group".id`

	services.DBCPsql.Raw(query, branchId, searchStr+"%").Scan(&sGroup)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   sGroup,
	})
}

func GroupDetail(ctx *iris.Context) {

	id := ctx.Get("id")
	groupBorrower := []GroupAgentBorrower{}

	query := "SELECT \"group\".\"id\",\"r_group_branch\".\"branchId\", \"group\".\"name\" as \"name\", \"group\".\"lat\" as \"lat\"," +
		"\"group\".\"lng\" as \"lng\", \"group\".\"scheduleDay\" as \"scheduleDay\", " +
		"\"group\".\"scheduleTime\" as \"scheduleTime\", \"group\".\"name\", cif.\"name\" as \"borrowerName\"," +
		"agent.fullname as \"agentName\",agent.id as \"agentId\",borrower.id as \"borrowerId\" "
	query += "FROM \"group\" "
	query += "LEFT JOIN r_group_borrower rgb ON rgb.\"groupId\" = \"group\".\"id\" "
	query += "LEFT JOIN borrower ON borrower.\"id\" = rgb.\"borrowerId\" "
	query += "LEFT JOIN r_cif_borrower rcb ON rcb.\"borrowerId\" = borrower.\"id\" "
	query += "lEFT join r_group_agent on r_group_agent.\"groupId\"=\"group\".id "
	query += "lEFT join agent on agent.id =  r_group_agent.\"agentId\" "
	query += "lEFT join r_group_branch on r_group_branch.\"groupId\" =  \"group\".id "
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

func Create(ctx *iris.Context) {

	type Payload struct {
		ID           uint64  `json:"_id"`
		Name         string  `json:"name"`
		ScheduleDay  string  `json:"scheduleDay"`
		ScheduleTime string  `json:"scheduleTime"`
		Lat          float64 `json:"lat"`
		Lng          float64 `json:"lng"`
		Agent        uint64  `json:"agentId"`
		Branch       uint64  `json:"branchId"`
	}

	m := Payload{}
	err := ctx.ReadJSON(&m)

	g := Group{}
	g.Name = m.Name
	g.ScheduleDay = m.ScheduleDay
	g.ScheduleTime = m.ScheduleTime
	g.Lat = m.Lat
	g.Lng = m.Lng

	if err != nil {
		panic(err)
	} else {
		services.DBCPsql.Create(&g)

		rga := r.RGroupAgent{}
		rga.GroupId = g.ID
		rga.AgentId = m.Agent

		if err := services.DBCPsql.Create(&rga).Error; err != nil {
			panic(err)
		}

		rgb := r.RGroupBranch{}
		rgb.GroupId = g.ID
		rgb.BranchId = m.Branch

		services.DBCPsql.Create(&rgb)

		rgborrower := r.RGroupBorrower{GroupId: g.ID, BorrowerId: 0}
		services.DBCPsql.Create(&rgborrower)

	}
	ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": m})

}

func Update(ctx *iris.Context) {
	groupId := ctx.Get("id")
	type Payload struct {
		ID           uint64  `json:"_id"`
		Name         string  `json:"name"`
		ScheduleDay  string  `json:"scheduleDay"`
		ScheduleTime string  `json:"scheduleTime"`
		Lat          float64 `json:"lat"`
		Lng          float64 `json:"lng"`
		Agent        uint64  `json:"agentId"`
	}

	m := Payload{}
	err := ctx.ReadJSON(&m)

	g := Group{}
	g.Name = m.Name
	g.ScheduleDay = m.ScheduleDay
	g.ScheduleTime = m.ScheduleTime
	g.Lat = m.Lat
	g.Lng = m.Lng

	query := `update "group" set "name" = ?, "scheduleDay" = ?, "scheduleTime" = ?, "lat" = ?, "lng" = ? where "group"."id" = ?`
	if err != nil {
		panic(err)
	} else {
		services.DBCPsql.Raw(query, g.Name, g.ScheduleDay, g.ScheduleTime, g.Lat, g.Lng, groupId).Scan(&g)

		rga := r.RGroupAgent{}
		rga.AgentId = m.Agent

		query_r_group_agent := `update "r_group_agent" set "agentId" = ? where "r_group_agent"."groupId" = ?`
		services.DBCPsql.Raw(query_r_group_agent, m.Agent, groupId).Scan(&rga)
	}

}

func UpdateGroupBorrower(ctx *iris.Context) {
	groupId := ctx.Get("id")

	type Payload struct {
		BorrowerId uint64 `json:"borrowerId"`
	}
	m := Payload{}
	err := ctx.ReadJSON(&m)

	if err != nil {
		panic(err)
	} else {
		r := r.RGroupBorrower{}
		r.BorrowerId = m.BorrowerId

		query := `update "r_group_borrower" set "borrowerId" = ? where "r_group_borrower"."groupId" = ?`
		services.DBCPsql.Raw(query, r.BorrowerId, groupId).Scan(&r)
	}
}

// GetGroupByBranchID is a method to get group by branch ID
func GetGroupByBranchID(ctx *iris.Context) {
	query := `SELECT "group".* FROM "group"
	JOIN r_group_branch ON r_group_branch."groupId" = "group".id
	WHERE "r_group_branch"."branchId" = ?`

	groupSchema := []Group{}

	services.DBCPsql.Raw(query, ctx.Param("branch_id")).Scan(&groupSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   groupSchema,
	})
}
