package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

var (
	DefaultApiPath   string
	Domain           string
	Port             string
	LogMode          bool
	PsqlHostAddress  string
	MysqlHostAddress string
)

type Config struct {
	Psql  DbConfig `json:"psql"`
	Mysql DbConfig `json:"mysql"`
}

type DbConfig struct {
	Db       string `json:"db"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SslMode  string `json:"ssl_mode"`
}

func init() {
	// Command flags
	defaultApiPath := flag.String("api-path", "/api/v1", "Default API path")
	domain := flag.String("domain", "default", "Service domain")
	port := flag.String("port", ":8080", "Port number")
	configPath := flag.String("config", "./config.json", "Config file path")
	logMode := flag.Bool("log-mode", false, "Log mode")

	// Parse all command flags
	flag.Parse()

	// Set service domain and application port
	DefaultApiPath = *defaultApiPath
	Domain = *domain
	Port = *port
	LogMode = *logMode

	// Setup config
	configFile, e := ioutil.ReadFile(*configPath)
	if e != nil {
		panic(e)
	}

	var c Config
	json.Unmarshal(configFile, &c)

	// Postgresql Host Address
	PsqlHostAddress = "host=" + c.Psql.Host + " port=" + c.Psql.Port + " user=" + c.Psql.Username + " dbname=" + c.Psql.Db + " sslmode=" + c.Psql.SslMode
	if c.Psql.Password != "" {
		PsqlHostAddress = PsqlHostAddress + " password=" + c.Psql.Password
	}

	// Mysql Host Address
	MysqlHostAddress = "host=" + c.Mysql.Host + " port=" + c.Mysql.Port + " user=" + c.Mysql.Username + " dbname=" + c.Mysql.Db + " sslmode=" + c.Mysql.SslMode
	if c.Mysql.Password != "" {
		MysqlHostAddress = MysqlHostAddress + " password=" + c.Mysql.Password
	}
}
