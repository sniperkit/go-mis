package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	P2pPath                   string
	WhiteList                 string
	GoCasPath                 string
	EnableEmergencyLoan       bool
	GoBankingPath             string
	GoLogPath                 string
	GoBorrowerPath            string
	GoLoanPath                string
	GoWithdrawalPath          string
	FlagServerPath            string
	Configuration             Config
)

type Config struct {
	Psql                []DbConfig  `json:"psql"`
	Mysql               DbConfig    `json:"mysql"`
	Redis               RedisConfig `json:"redis"`
	UploaderPath        string      `json:"uploaderPath"`
	P2pPath             string      `json:"p2pPath"`
	GoCasPath           string      `json:"goCasPath"`
	WhiteList           string      `json:"whiteList"`
	SignString          string      `json:"signString"`
	ApiVersion          string      `json:"apiVersion"`
	EnableEmergencyLoan bool        `json:"enableEmergencyLoan"`
	GoWithdrawalPath    string      `json:"goWithdrawalPath"`
	GoBankingPath       string      `json:"goBankingPath"`
	GoLogPath           string      `json:"goLogPath"`
	GoLoanPath          string      `json:"goLoanPath"`
	GoBorrowerPath      string      `json:"goBorrowerPath"`
	FlagServerPath      string      `json:"flagServerPath"`
}

type DbConfig struct {
	Db       string `json:"db"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SslMode  string `json:"ssl_mode"`
}

// RedisConfig - redis configuration
type RedisConfig struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Db       int    `json:"db"`
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
	// var c Config
	err := json.Unmarshal(configFile, &Configuration)
	if err != nil {
		log.Println("[ERROR] Error when reading configuration file")
		panic(err)
	}
	Configuration.ApiVersion = "2.2.0"
	fmt.Println("Version:", Configuration.ApiVersion)
	fmt.Println("------------------")
	Version = Configuration.ApiVersion

	UploaderApiPath = Configuration.UploaderPath
	P2pPath = Configuration.P2pPath
	GoCasApiPath = Configuration.GoCasPath
	GoLogPath = Configuration.GoLogPath
	SignStringKey = Configuration.SignString
	GoBankingPath = Configuration.GoBankingPath
	GoBorrowerPath = Configuration.GoBorrowerPath
	GoLoanPath = Configuration.GoLoanPath
	GoWithdrawalPath = Configuration.GoWithdrawalPath
	EnableEmergencyLoan = Configuration.EnableEmergencyLoan
	WhiteList = Configuration.WhiteList
	FlagServerPath = Configuration.FlagServerPath

	// Postgresql Host Address
	PsqlHostAddressMisAmartha = "host=" + Configuration.Psql[0].Host + " port=" + Configuration.Psql[0].Port + " user=" + Configuration.Psql[0].Username + " dbname=" + Configuration.Psql[0].Db + " sslmode=" + Configuration.Psql[0].SslMode + " fallback_application_name=go-mis"
	if Configuration.Psql[0].Password != "" {
		PsqlHostAddressMisAmartha += " password=" + Configuration.Psql[0].Password
	}

	PsqlHostAddressSurvey = "host=" + Configuration.Psql[1].Host + " port=" + Configuration.Psql[1].Port + " user=" + Configuration.Psql[1].Username + " dbname=" + Configuration.Psql[1].Db + " sslmode=" + Configuration.Psql[1].SslMode
	if Configuration.Psql[1].Password != "" {
		PsqlHostAddressSurvey += " password=" + Configuration.Psql[1].Password
	}

	// Mysql Host Address
	MysqlHostAddress = "host=" + Configuration.Mysql.Host + " port=" + Configuration.Mysql.Port + " user=" + Configuration.Mysql.Username + " dbname=" + Configuration.Mysql.Db + " sslmode=" + Configuration.Mysql.SslMode
	if Configuration.Mysql.Password != "" {
		MysqlHostAddress = MysqlHostAddress + " password=" + Configuration.Mysql.Password
	}
}
