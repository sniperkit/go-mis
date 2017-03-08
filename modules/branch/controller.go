package branch

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Branch{})
	services.BaseCrudInit(Branch{}, []Branch{})
}

// FetchAll - fetchAll branchs data
/** Todo error handling */
func FetchAll(ctx *iris.Context) {
	// branch := getBranchWithoutManager()
	// branchManager := getBranchManager()
	// branch = combineBranchManager(branch, branchManager)

	query := "SELECT branch.id, area.\"name\" AS \"area\", branch.\"name\" AS \"name\", branch.city, branch.province, user_mis.fullname AS \"manager\", \"role\".\"name\" AS \"role\" "
	query += "FROM branch "
	query += "LEFT JOIN r_area_branch ON r_area_branch.\"branchId\" = branch.id "
	query += "LEFT JOIN area ON area.id = r_area_branch.\"areaId\" "
	query += "LEFT JOIN r_branch_user_mis ON r_branch_user_mis.\"branchId\" = branch.id "
	query += "LEFT JOIN user_mis ON user_mis.id = r_branch_user_mis.\"branchId\" "
	query += "LEFT JOIN r_user_mis_role ON r_user_mis_role.\"userMisId\" = user_mis.id "
	query += "LEFT JOIN \"role\" ON \"role\".id = r_user_mis_role.\"roleId\" "
	query += "WHERE (\"role\".\"name\" ~* 'branch manager' OR \"role\".id IS NULL) AND branch.\"deletedAt\" IS NULL "

	branch := []BranchManagerArea{}
	services.DBCPsql.Raw(query).Scan(&branch)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   branch,
	})
}

func getBranchWithoutManager() []BranchManagerArea {
	query := "SELECT branch.\"id\", branch.\"name\", branch.\"city\", branch.\"province\", area.\"name\" as \"area\" "
	query += "FROM branch "
	query += "LEFT JOIN r_area_branch ON r_area_branch.\"branchId\" = branch.\"id\" "
	query += "LEFT JOIN area ON r_area_branch.\"areaId\" = area.\"id\" "
	query += "WHERE branch.\"deletedAt\" IS NULL"

	result := []BranchManagerArea{}
	services.DBCPsql.Raw(query).Find(&result)
	return result
}

func getBranchManager() []BranchManager {
	query := "SELECT \"branchId\",\"fullname\" "
	query += "FROM r_branch_user_mis rbum "
	query += "JOIN user_mis um on um.id = rbum.\"userMisId\" "
	query += "LEFT JOIN r_user_mis_role rumr on rbum.\"userMisId\" = rumr.\"userMisId\" "
	//query += "WHERE roleId = 4 "

	result := []BranchManager{}
	services.DBCPsql.Raw(query).Find(&result)
	return result
}

func combineBranchManager(branch []BranchManagerArea, branchManager []BranchManager) []BranchManagerArea {
	for i := 0; i < len(branch); i++ {
		for j := 0; j < len(branchManager); j++ {
			if branch[i].ID == branchManager[j].BranchId {
				branch[i].Manager = branchManager[j].Fullname
			}
		}
	}
	return branch
}

// GetByID branch by id
func GetByID(ctx *iris.Context) {
	bracnh := Branch{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).First(&bracnh)
	if bracnh == (Branch{}) {
		ctx.JSON(iris.StatusOK, iris.Map{"data": iris.Map{}})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": bracnh})
	}
}
