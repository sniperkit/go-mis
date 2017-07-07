package userMis

import (
	"time"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&UserMis{})
	services.BaseCrudInit(UserMis{}, []UserMis{})
}

func CreateUserMis(ctx *iris.Context){
	type Payload struct {
		ID          uint64     `json:"_id"`
		Fullname    string     `json:"fullname"`
		Username    string     `json:"username"`
		Password    string     `json:"password"`
		PhoneNo     string     `json:"phoneNo"`
		PicUrl      string     `json:"picUrl"`
		Role    	 	uint64     `json:"role"`
		Area  	    uint64     `json:"area"`
		Branch      uint64  `json:"branch"`
	}
	m := Payload{}
	err := ctx.ReadJSON(&m)
	userMis := UserMis{}


	userMis.ID = m.ID
	userMis.Fullname = m.Fullname;
	userMis.Username = m.Username;
	userMis.Password = m.Password;
	userMis.PhoneNo = m.PhoneNo;
	userMis.PicUrl = m.PicUrl;


	if err != nil{
		panic(err)
	}else{
		services.DBCPsql.Create(&userMis);

		rur := r.RUserMisRole{}
		rur.UserMisId = userMis.ID;
		rur.RoleId = m.Role;
		services.DBCPsql.Create(&rur);

		rbu := r.RBranchUserMis{}
		rbu.UserMisId = userMis.ID;
		rbu.BranchId = m.Branch;
		services.DBCPsql.Create(&rbu);

		rau := r.RAreaUserMis{}
		rau.UserMisId = userMis.ID;
		rau.AreaId = m.Area;
		services.DBCPsql.Create(&rau);

	}
	ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": m})

}

func UpdateUserBranch(ctx *iris.Context) {
	userObj := ctx.Get("USER_MIS").(UserMis)
	branchId := ctx.Get("branch_id")

	userMisBranch := r.RBranchUserMis{}
	query := "update r_branch_user_mis set \"branchId\" = ? where \"userMisId\" = ?"
	services.DBCPsql.Raw(query, branchId, userObj.ID).Scan(&userMisBranch)
}

func GetUserMisById(ctx *iris.Context){
	userMis := UserMisAreaBranchRole{}
	id := ctx.Get("id")

	query := `SELECT user_mis."id" AS "userMisId", user_mis."phoneNo", user_mis."_password", user_mis."_username", user_mis."picUrl", user_mis."fullname", user_mis."isSuspended", role."name" AS role,"role"."id" as "roleId", branch."id" as "branchId", area."id" as "areaId" , branch."name" AS "branch", area."name" AS "area" FROM user_mis
		LEFT JOIN r_branch_user_mis ON r_branch_user_mis."userMisId" = user_mis."id" 
		LEFT JOIN branch ON branch."id" = r_branch_user_mis."branchId" 
		LEFT JOIN r_area_branch ON r_area_branch."branchId" = branch."id" 
		LEFT JOIN area ON area."id" = r_area_branch."areaId" 
		LEFT JOIN r_user_mis_role ON r_user_mis_role."userMisId" = user_mis."id" 
		LEFT JOIN role ON role."id" = r_user_mis_role."roleId" 
		WHERE user_mis."deletedAt" IS NULL AND user_mis."id" = ?`
	services.DBCPsql.Raw(query, id).Find(&userMis)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   userMis,
	})


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

func DeleteUserMis(ctx *iris.Context) {
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
