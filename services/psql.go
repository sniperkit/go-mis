package services

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"bitbucket.org/go-mis/config"
)

var DBCPsql *gorm.DB

func init() {
	var err error
	con, err := gorm.Open("postgres", config.PsqlHostAddress)

	if err != nil {
		fmt.Println("[ERROR] Failed to connect to postgres. Config=" + config.PsqlHostAddress)
		return
	}

	con.LogMode(config.LogMode)
	con.SingularTable(true)
	con.Exec("CREATE EXTENSION postgis")
	con.Exec("CREATE EXTENSION postgis_topology")

	DBCPsql = con
	fmt.Println("[INFO] Connected to PSQL. Config => " + config.PsqlHostAddress)
}
