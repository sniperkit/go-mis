package services

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"bitbucket.org/go-mis/config"
)

var DBCPsql *gorm.DB
var DBCPsqlSurvey *gorm.DB

func init() {
	var err error
	con, err := gorm.Open("postgres", config.PsqlHostAddressMisAmartha)

	if err != nil {
		fmt.Println("[ERROR] Failed to connect to postgres. Config=" + config.PsqlHostAddressMisAmartha)
		return
	}

	con.LogMode(config.LogMode)
	con.SingularTable(true)
	con.Exec("CREATE EXTENSION postgis")
	con.Exec("CREATE EXTENSION postgis_topology")
	con.DB().SetMaxIdleConns(10)
	con.DB().SetMaxOpenConns(40)

	DBCPsql = con
	fmt.Println("[INFO] Connected to PSQL. Config => " + config.PsqlHostAddressMisAmartha)

	var errSurvey error
	conSurvey, errSurvey := gorm.Open("postgres", config.PsqlHostAddressSurvey)

	if errSurvey != nil {
		fmt.Println("[ERROR] Failed to connect to postgres. Config=" + config.PsqlHostAddressSurvey)
		return
	}

	conSurvey.LogMode(config.LogMode)
	conSurvey.SingularTable(true)
	conSurvey.DB().SetMaxIdleConns(10)
	conSurvey.DB().SetMaxOpenConns(85)

	DBCPsqlSurvey = conSurvey
	fmt.Println("[INFO] Connected to PSQL. Config => " + config.PsqlHostAddressSurvey)
}
