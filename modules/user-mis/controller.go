package userMis

import (
	"time"
	"net/http"
	"encoding/json"
	"bytes"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
	"github.com/jinzhu/gorm"
	"fmt"
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
		Role    		uint64     `json:"role"`
		Area  	    uint64     `json:"area"`
		Branch      uint64     `json:"branch"`
	}

	type Cas struct{
    Username string `json:"username"`
    Password string `json:"password"`
		UserType string `json:"userType"`
	}

	m := Payload{}
	err := ctx.ReadJSON(&m)

	if err != nil{
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}else{
		userMis := UserMis{}

		userMis.ID = m.ID
		userMis.Fullname = m.Fullname;
		userMis.Username = m.Username;
		userMis.Password = m.Password;
		userMis.PhoneNo = m.PhoneNo;
		userMis.PicUrl = m.PicUrl;

		/***
		Function Register to Go-Cas
		author 	: @primayudantra
		date		: 24 July 2017
		***/

		u := Cas{Username: userMis.Username, Password: userMis.Password, UserType: "MIS"}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(u)
    res, _ := http.Post("http://localhost:4500/api/v1/register", "application/json; charset=utf-8", b)

		if res.Status == "200 OK"{
			db:=services.DBCPsql.Begin()
			if err:=db.Create(&userMis).Error;err!=nil{
				processErrorAndRollback(ctx, db, err, "Create User")
				return
			};

			rur := r.RUserMisRole{}
			rur.UserMisId = userMis.ID;
			rur.RoleId = m.Role;
			if err:=db.Create(&rur).Error;err!=nil{
				processErrorAndRollback(ctx, db, err, "Create User")
				return
			};

			rbu := r.RBranchUserMis{}
			rbu.UserMisId = userMis.ID;
			rbu.BranchId = m.Branch;
			if err:=db.Create(&rbu).Error;err!=nil{
				processErrorAndRollback(ctx, db, err, "Create User")
				return
			};

			rau := r.RAreaUserMis{}
			rau.UserMisId = userMis.ID;
			rau.AreaId = m.Area;
			if err:=db.Create(&rau).Error;err!=nil{
				processErrorAndRollback(ctx, db, err, "Create User")
				return
			};
			db.Commit()
		}else{
			ctx.JSON(iris.StatusUnauthorized, iris.Map{
				"status":  "error",
				"message": "Error Create User in Go-CAS",
			})
			return
		}
	}
	ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": m})

}

func UpdateUserMisById(ctx *iris.Context){
	type Payload struct {
		ID			uint64		`json:"id"`
		Fullname	string		`json:"fullname"`
		Username	string		`json:"_username"`
		PhoneNo		string		`json:"phoneNo"`
		PicUrl		string		`json:"picUrl"`
		Role		uint64		`json:"role"`
		Area		uint64		`json:"area"`
		Branch		uint64		`json:"branch"`
	}

	m := Payload{}
	err := ctx.ReadJSON(&m)


	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	} else {
		userMis := UserMis{}
		userMis.ID = m.ID
		userMis.Fullname = m.Fullname;
		userMis.Username = m.Username;
		userMis.PhoneNo = m.PhoneNo;
		userMis.PicUrl = m.PicUrl;

		db:=services.DBCPsql.Begin()
		// Update User
		userQuery := `UPDATE user_mis SET "fullname" = ?, "_username" = ?, "phoneNo" = ?, "picUrl" = ? WHERE "id" = ?`
		if err:=db.Exec(userQuery, userMis.Fullname, userMis.Username,userMis.PhoneNo, userMis.PicUrl, userMis.ID).Error;err!=nil {
			fmt.Println("Error",err)
			processErrorAndRollback(ctx, db, err, "Update user")
			return
		};

		// Update Role
		updateRole := r.RUserMisRole{}
		updateRole.UserMisId = userMis.ID;
		updateRole.RoleId = m.Role;
		roleQuery := `UPDATE r_user_mis_role SET "roleId" = ? where "userMisId" = ?`
		if err=db.Exec(roleQuery, updateRole.RoleId, updateRole.UserMisId).Error;err!=nil {
			fmt.Println("Error",err)
			processErrorAndRollback(ctx, db, err, "Update user")
			return
		}

		// Update BranchId
		updateBranch := r.RBranchUserMis{}
		updateBranch.UserMisId = userMis.ID;
		updateBranch.BranchId = m.Branch;
		branchQuery := `UPDATE r_branch_user_mis SET "branchId" = ? where "userMisId" = ?`
		if err:=db.Exec(branchQuery, updateBranch.BranchId, updateBranch.UserMisId).Error;err!=nil {
			fmt.Println("Error",err)
			processErrorAndRollback(ctx, db, err, "Update user")
			return
		}

		// Update Area
		updateArea := r.RAreaUserMis{}
		updateArea.UserMisId = userMis.ID;
		updateArea.AreaId = m.Area;
		areaQuery := `UPDATE r_area_user_mis SET "areaId" = ? where "userMisId" = ?`
		if err:=db.Exec(areaQuery, updateArea.AreaId, updateArea.UserMisId).Error;err!=nil {
			fmt.Println("Error",err)
			processErrorAndRollback(ctx, db, err, "Update user")
			return
		}
		db.Commit()
	}

	ctx.JSON(iris.StatusOK, iris.Map {
		"status": "success",
		"data": m })

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

func processErrorAndRollback(ctx *iris.Context, db *gorm.DB, err error, process string) {
	db.Rollback()
	ctx.JSON(iris.StatusInternalServerError, iris.Map{"status": "error", "error": "Error on " + process + " " + err.Error()})
}
