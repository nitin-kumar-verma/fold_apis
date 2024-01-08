package main

import (
	"fold/internal/config"
	"fold/internal/routes"
	"fold/internal/startup"
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	config.LoadConfig()
	startup.InitServer()
	app := routes.SetupRoutes()
	err := app.Listen(":" + config.GlobalConfig.Port)
	if err != nil {
		log.Fatalf("Error starting server :%v", err)
	}
}
