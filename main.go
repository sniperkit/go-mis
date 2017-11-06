package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/iris-contrib/middleware.v4/cors"
	"gopkg.in/iris-contrib/middleware.v4/logger"
	"gopkg.in/iris-contrib/middleware.v4/recovery"
	"gopkg.in/kataras/iris.v4"

	"regexp"

	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/routes"
)

func main() {

	// Initialize recovery
	iris.Use(recovery.New())

	// Initialize logger
	iris.Use(logger.New())
	iris.Use(NewHttpLog())

	// Check environment, if `dev` then let the CORS to `*`
	if config.Env == "dev" || config.Env == "development" {
		crs := cors.New(cors.Options{})
		iris.Use(crs)
	}

	// Initialize routes
	routes.Init()

	// Initialize custom routes
	routes.InitCustomApi()
	// Start app
	iris.Listen(config.Port)
}

type customHTTPLoggerMiddleware struct {
}

// NewHttpLog  - instance customHTTPLoggerMiddleware
func NewHttpLog() iris.HandlerFunc {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "/var/log"
	}
	path := pwd + "/httplog"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}

	l := &customHTTPLoggerMiddleware{}
	return l.logCustom(path)
}
func (l *customHTTPLoggerMiddleware) logCustom(path string) iris.HandlerFunc {
	return func(ctx *iris.Context) {
		startTime := time.Now()
		startDate := startTime.Format("2006-01-02 15:04:05")
		curDate := startTime.Format("2006-01-02")
		logName := curDate + ".log"
		method := ctx.MethodString()
		requestURI := string(ctx.Request.RequestURI())
		requestBody := string(ctx.Request.Body())
		if len(requestBody) == 0 {
			requestBody = "{}"
		}

		ctx.Next()
		//no time.Since in order to format it well after
		endTime := time.Now()

		endDate := endTime.Format("01/02 - 15:04:05")
		status := strconv.Itoa(ctx.Response.StatusCode())
		responseBody := string(ctx.Response.Body())
		latency := endTime.Sub(startTime)

		requestLogString := fmt.Sprintf("#REQUEST %s %4v %s %s \n%s\n", startDate, latency, method, requestURI, requestBody)
		responseLogString := fmt.Sprintf("#RESPONSE %s %v %s\n\n", endDate, status, responseBody)
		logString := requestLogString + ", " + responseLogString

		//remove newline
		regex := regexp.MustCompile(`\r?\n`)
		logString = regex.ReplaceAllString(logString, "") + "\n"
		l.createAndWrite(path+"/"+logName, logString)
	}
}

func (l *customHTTPLoggerMiddleware) createAndWrite(path string, log string) {
	l.createFile(path)
	l.writeFile(path, log)
}

func (l *customHTTPLoggerMiddleware) createFile(path string) {
	fmt.Print(path)
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			fmt.Println(err.Error())

		}
		defer file.Close()
	}
}

func (l *customHTTPLoggerMiddleware) writeFile(path string, txt string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err.Error())

	}
	defer file.Close()

	// write some text to file
	_, err = file.WriteString(txt)
	if err != nil {
		fmt.Println(err.Error())

	}

	// save changes
	err = file.Sync()
	if err != nil {
		fmt.Println(err.Error())

	}
}
