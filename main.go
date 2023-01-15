package main

import (
	"birdnest/app"
	"birdnest/docs"
	"birdnest/logger"
)

// @title          Birdnest API
// @version        1.0
// @description    This is a birdnest web application.
// @BasePath  /api

func main() {
	cnf := app.ReadConfig()
	docs.SwaggerInfo.Host = cnf.SwagHost
	APP, err := app.NewApp(cnf)
	if err != nil {
		logger.AppLogger.Info("failure to create app structure:", err)
		return
	}
	app.StartRouter(APP)
}
