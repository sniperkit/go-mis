package agent

import (
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
	"github.com/jinzhu/gorm"
	"bitbucket.org/go-mis/modules/r"
	"fmt"
	"strconv"
	"bytes"
	"encoding/json"
	"net/http"
	"bitbucket.org/go-mis/config"
)

type Cas struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType string `json:"userType"`
}

func Init() {
	services.DBCPsql.AutoMigrate(&Agent{})
	services.BaseCrudInit(Agent{}, []Agent{})
}

func GetAllAgentByBranchID(ctx *iris.Context) {
	branchID := ctx.Get("BRANCH_ID")
	GetAgent(ctx, branchID.(uint64))
}

func GetAgent(ctx *iris.Context, branchID uint64) {
	agentSchema := []Agent{}
	query := ""

	query += "SELECT agent.id, agent.\"picUrl\", agent.\"username\", agent.fullname, agent.address,  "
	query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = '0' AND kecamatan = '0' AND kelurahan = '0' LIMIT 1) AS \"province\", "
	query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = agent.city AND kecamatan = '0' AND kelurahan = '0' LIMIT 1) AS \"city\", "
	query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = agent.city AND kecamatan = agent.kecamatan AND kelurahan = '0' LIMIT 1) AS \"kecamatan\", "
	query += "(SELECT \"name\" FROM inf_location WHERE province = agent.province AND city = agent.city AND kecamatan = agent.kecamatan AND kelurahan = agent.kelurahan LIMIT 1) AS \"kelurahan\" "
	query += "FROM agent "
	query += "INNER JOIN r_branch_agent ON r_branch_agent.\"agentId\" = agent.id "
	query += "WHERE r_branch_agent.\"branchId\" = ? AND agent.\"deletedAt\" IS NULL"

	services.DBCPsql.Raw(query, branchID).Scan(&agentSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   agentSchema,
	})
}

func GetAllAgent(ctx *iris.Context) {

	id := ctx.Get("id")
	branchId, _ := strconv.ParseUint(id.(string), 0, 64)
	GetAgent(ctx, branchId)

}

func GetAgentById(ctx *iris.Context) {
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

func CreateAgent(ctx *iris.Context) {

	type Payload struct {
		Username        string           `json:"username"`
		Password        string           `json:"password"`
		Fullname        string           `json:"fullname"`
		BankName        string           `json:"bankName"`
		BankAccountName string           `json:"bankAccountName"`
		BankAccountNo   string           `json:"bankAccountNo"`
		PicUrl          string           `json:"picUrl"`
		PhoneNo         string           `json:"phoneNo"`
		Address         string           `json:"address"`
		Kelurahan       string           `json:"kelurahan"`
		Kecamatan       string           `json:"kecamatan"`
		City            string           `json:"city"`
		Province        string           `json:"province"`
		Nationality     string           `json:"nationality"`
		Lat             float64          `json:"lat"`
		Lng             float64          `json:"lng"`
		Branch          uint64             `json:"branchId"`
	}

	m := Payload{}
	err := ctx.ReadJSON(&m)
	//Register to Go-CAS
	u := Cas{Username: m.Username, Password: m.Password, UserType: "APP"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	res, _ := http.Post(config.GoCasApiPath+"/api/v1/register", "application/json; charset=utf-8", b)
	if res.StatusCode != 200 {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Error Create User in Go-CAS",
		})
		return
	}
	//Create Agent to Postgre
	a := Agent{}
	a.Username = m.Username;
	a.Fullname = m.Fullname;
	a.BankName = m.BankName;
	a.BankAccountName = m.BankAccountName;
	a.BankAccountNo = m.BankAccountNo;
	a.PicUrl = m.PicUrl;
	a.PhoneNo = m.PhoneNo;
	a.Address = m.Address;
	a.Kelurahan = m.Kelurahan;
	a.Kecamatan = m.Kecamatan;
	a.City = m.City;
	a.Province = m.Province;
	a.Nationality = m.Nationality;
	a.Lat = m.Lat;
	a.Lng = m.Lng;
	a.Password = m.Password
	fmt.Println("Password", m.Password)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	} else {
		db := services.DBCPsql.Begin()
		if err := db.Create(&a).Error; err != nil {
			processErrorAndRollback(ctx, db, err, "Create agent")
			return

		}
		rba := r.RBranchAgent{}
		rba.AgentId = a.ID;
		rba.BranchId = m.Branch;
		if err := db.Create(&rba).Error; err != nil {
			processErrorAndRollback(ctx, db, err, "Create agent")
			return
		};
		db.Commit()
	}
	ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": m})
}

func UpdateAgent(ctx *iris.Context) {
	agentId := ctx.Get("id")
	type Payload struct {
		Username        string           `json:"username"`
		Password        string           `json:"password"`
		Fullname        string           `json:"fullname"`
		BankName        string           `json:"bankName"`
		BankAccountName string           `json:"bankAccountName"`
		BankAccountNo   string           `json:"bankAccountNo"`
		PicUrl          string           `json:"picUrl"`
		PhoneNo         string           `json:"phoneNo"`
		Address         string           `json:"address"`
		Kelurahan       string           `json:"kelurahan"`
		Kecamatan       string           `json:"kecamatan"`
		City            string           `json:"city"`
		Province        string           `json:"province"`
		Nationality     string           `json:"nationality"`
		Lat             float64          `json:"lat"`
		Lng             float64          `json:"lng"`
		Branch          uint64             `json:"branchId"`
	}

	m := Payload{}
	err := ctx.ReadJSON(&m)
	//Update Agent to Postgre
	a := Agent{}
	a.Username = m.Username;
	a.Password = m.Password
	a.Fullname = m.Fullname;
	a.BankName = m.BankName;
	a.BankAccountName = m.BankAccountName;
	a.BankAccountNo = m.BankAccountNo;
	a.PicUrl = m.PicUrl;
	a.PhoneNo = m.PhoneNo;
	a.Address = m.Address;
	a.Kelurahan = m.Kelurahan;
	a.Kecamatan = m.Kecamatan;
	a.City = m.City;
	a.Province = m.Province;
	a.Nationality = m.Nationality;
	a.Lat = m.Lat;
	a.Lng = m.Lng;
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	} else {
		db := services.DBCPsql.Begin()
		if err := db.Table("agent").Where(" \"id\" = ?", agentId).Update(&a).Error; err != nil {
			processErrorAndRollback(ctx, db, err, "Update agent")
			return
		}
		db.Commit()
	}
	ctx.JSON(iris.StatusOK, iris.Map{"status": "success", "data": a})

}

func UpdateAgentPasswordById(ctx *iris.Context) {
	agentId := ctx.Get("id")
	fmt.Println(agentId)

	type Payload struct {
		Username string        `json:"username"`
		Password string        `json:"password"`
	}

	type Cas struct {
		Username string `json:"username"`
		Password string `json:"password"`
		UserType string `json:"userType"`
	}

	m := Payload{}
	err := ctx.ReadJSON(&m)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}

	agent := Agent{}
	agent.Username = m.Username
	agent.Password = m.Password;
	u := Cas{Username: agent.Username, Password: agent.Password, UserType: "APP"}
	fmt.Println(u)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	res, _ := http.Post(config.GoCasApiPath+"/api/v1/update-password", "application/json; charset=utf-8", b)
	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)
	if res.StatusCode != 200 {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Error update password in go-cas",
		})
		return
	}
	agent.BeforeUpdate()
	db := services.DBCPsql.Begin()
	// Update User in PSQL
	userQuery := `UPDATE agent SET "password" = ? WHERE "id" = ?`
	db.Exec(userQuery, agent.Password, agentId)
	db.Commit();

}

func processErrorAndRollback(ctx *iris.Context, db *gorm.DB, err error, process string) {
	db.Rollback()
	ctx.JSON(iris.StatusInternalServerError, iris.Map{"status": "error", "error": "Error on " + process + " " + err.Error()})
}
