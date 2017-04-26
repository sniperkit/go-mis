package userMis

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
	"bitbucket.org/go-mis/modules/r"
	"time"
)

func Init() {
	services.DBCPsql.AutoMigrate(&UserMis{})
	services.BaseCrudInit(UserMis{}, []UserMis{})
}

func FetchUserMisAreaBranchRole(ctx *iris.Context) {
	arrUserMisAreaBranchRole := []UserMisAreaBranchRole{}

	query := "SELECT user_mis.\"id\" AS \"userMisId\", user_mis.\"picUrl\", user_mis.\"fullname\", user_mis.\"isSuspended\", role.\"name\" AS \"role\", branch.\"name\" AS \"branch\", area.\"name\" AS \"area\" "
	query += "FROM user_mis "
	query += "LEFT JOIN r_branch_user_mis ON r_branch_user_mis.\"userMisId\" = user_mis.\"id\" "
	query += "LEFT JOIN branch ON branch.\"id\" = r_branch_user_mis.\"branchId\" "
	query += "LEFT JOIN r_area_branch ON r_area_branch.\"branchId\" = branch.\"id\" "
	query += "LEFT JOIN area ON area.\"id\" = r_area_branch.\"areaId\" "
	query += "LEFT JOIN r_user_mis_role ON r_user_mis_role.\"userMisId\" = user_mis.\"id\" "
	query += "LEFT JOIN role ON role.\"id\" = r_user_mis_role.\"roleId\" "
	query += "WHERE user_mis.\"deletedAt\" IS NULL "

	services.DBCPsql.Raw(query).Find(&arrUserMisAreaBranchRole)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   arrUserMisAreaBranchRole,
	})
}

func DeleteUserMis (ctx *iris.Context) {
	// delete user
	m := UserMis{}
	services.DBCPsql.Model(m).Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).UpdateColumn("deletedAt", time.Now())

	// delete relation to role, area, and branch
	var mRel []interface{}
	mRel = append(mRel, r.RUserMisRole{}, r.RBranchUserMis{}, r.RAreaUserMis{})
	for _, val := range mRel {
		services.DBCPsql.Model(val).Where("\"deletedAt\" IS NULL AND \"userMisId\" = ?", ctx.Param("id")).UpdateColumn("deletedAt", time.Now())
	}

	ctx.JSON(iris.StatusOK, iris.Map{"data": m})

}
