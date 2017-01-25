package main

import (
	"gopkg.in/iris-contrib/middleware.v4/cors"
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
