package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

var (
	ApiKey                    string
	DefaultApiPath            string
	Domain                    string
	Port                      string
	Env                       string
	LogMode                   bool
	PsqlHostAddressMisAmartha string
	PsqlHostAddressSurvey     string
	MysqlHostAddress          string
	UploaderApiPath           string
	GoCasPath			           string
	EnableEmergencyLoan		  bool
)

type Config struct {
	Psql         []DbConfig `json:"psql"`
	Mysql        DbConfig   `json:"mysql"`
	UploaderPath string     `json:"uploaderPath"`
	GoCasPath		 string     `json:"goCasPath"`
	SignString		 string     `json:"signString"`
	ApiVersion   string     `json:"apiVersion"`
	EnableEmergencyLoan bool `json:"enableEmergencyLoan"`
}

type DbConfig struct {
	Db       string `json:"db"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SslMode  string `json:"ssl_mode"`
}

var Version string
var GoCasApiPath string
var SignStringKey string

func init() {
	ApiKey = "$2a$06$20EpVmcNvVg0heEijxLEP.Aw0hhoC7kJyuGltJnYZMStuhOLwPB7W"

	// Command flags
	defaultApiPath := flag.String("api-path", "/api/v1", "Default API path")
	domain := flag.String("domain", "default", "Service domain")
	port := flag.String("port", ":8080", "Port number")
	configPath := flag.String("config", "./config.json", "Config file path")
	logMode := flag.Bool("log-mode", false, "Log mode")
	env := flag.String("env", "production", "Default environment")

	// Parse all command flags
	flag.Parse()

	// Set service domain and application port
	DefaultApiPath = *defaultApiPath
	Domain = *domain
	Port = *port
	LogMode = *logMode
	Env = *env

	fmt.Println("------------------")
	fmt.Println("Default API Path -", DefaultApiPath)
	fmt.Println("Application Port -", Port)
	fmt.Println("Enable log mode -", LogMode)
	fmt.Println("------------------")

	// Setup config
	configFile, e := ioutil.ReadFile(*configPath)
	if e != nil {
		panic(e)
	}

	var c Config
	json.Unmarshal(configFile, &c)
	c.ApiVersion = "2.1.0"
	fmt.Println("Version:", c.ApiVersion)
	fmt.Println("------------------")

	Version = c.ApiVersion

	UploaderApiPath = c.UploaderPath
	GoCasApiPath = c.GoCasPath
	SignStringKey = c.SignString



	EnableEmergencyLoan = c.EnableEmergencyLoan

	// Postgresql Host Address
	PsqlHostAddressMisAmartha = "host=" + c.Psql[0].Host + " port=" + c.Psql[0].Port + " user=" + c.Psql[0].Username + " dbname=" + c.Psql[0].Db + " sslmode=" + c.Psql[0].SslMode
	if c.Psql[0].Password != "" {
		PsqlHostAddressMisAmartha += " password=" + c.Psql[0].Password
	}

	PsqlHostAddressSurvey = "host=" + c.Psql[1].Host + " port=" + c.Psql[1].Port + " user=" + c.Psql[1].Username + " dbname=" + c.Psql[1].Db + " sslmode=" + c.Psql[1].SslMode
	if c.Psql[1].Password != "" {
		PsqlHostAddressSurvey += " password=" + c.Psql[1].Password
	}

	// Mysql Host Address
	MysqlHostAddress = "host=" + c.Mysql.Host + " port=" + c.Mysql.Port + " user=" + c.Mysql.Username + " dbname=" + c.Mysql.Db + " sslmode=" + c.Mysql.SslMode
	if c.Mysql.Password != "" {
		MysqlHostAddress = MysqlHostAddress + " password=" + c.Mysql.Password
	}
}
