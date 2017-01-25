package main

import (
	"gopkg.in/iris-contrib/middleware.v4/logger"
	"gopkg.in/iris-contrib/middleware.v4/recovery"
	"gopkg.in/kataras/iris.v4"

	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/routes"
)

func main() {

	// Initialize recovery
	iris.Use(recovery.New())

	// Initialize logger
	iris.Use(logger.New())

	// Initialize routes
	routes.Init()

	// Initialize custom routes
	routes.InitCustomApi()

	// Start app
	iris.Listen(config.Port)
}
