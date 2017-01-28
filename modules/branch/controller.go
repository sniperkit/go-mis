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
func FetchAll(ctx *iris.Context) {
	bracnhs := []BranchManager{}

	query := "SELECT branch.\"id\", branch.\"name\", branch.\"city\", branch.\"province\", user_mis.\"fullname\" as \"managerName\", area.\"name\" as \"areaName\" "
	query += "FROM branch "
	query += "JOIN r_branch_user_mis ON r_branch_user_mis.\"branchId\" = branch.\"id\" "
	query += "JOIN user_mis ON user_mis.\"id\" = r_branch_user_mis.\"userMisId\" "
	query += "JOIN r_user_mis_role ON r_user_mis_role.\"userMisId\" = user_mis.\"id\" "
	query += "JOIN role ON role.\"id\" = r_user_mis_role.\"roleId\" "
	query += "JOIN r_area_branch ON r_area_branch.\"branchId\" = branch.\"id\" "
	query += "JOIN area ON r_area_branch.\"areaId\" = area.\"id\" "
	query += "WHERE role.\"name\" = ?"

	services.DBCPsql.Raw(query, "branchmanager").Find(&bracnhs)
	ctx.JSON(iris.StatusOK, iris.Map{"data": bracnhs})
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
