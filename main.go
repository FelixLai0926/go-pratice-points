package main

import (
	"log"
	"os"
	"points/pkg/middleware"
	"points/pkg/module/config"
	"points/pkg/module/database"
	"points/pkg/module/logger"
	"points/pkg/router"
)

func main() {
	env := "example"
	if err := logger.InitLogger(env); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.SyncLogger()

	server := router.Setup()
	envFilePath, err := config.GetEnvPath(env)
	if err != nil {
		log.Fatalf("Error getting .env file path: %v", err)
	}
	if err := config.InitEnv(envFilePath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.ParseEnv[database.PostgresConfig]()
	if err != nil {
		log.Fatalf("Error transforming .env file to struct: %v", err)
	}

	gormdb, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Error initializing database:: %v", err)
	}

	sqlDB, err := gormdb.DB()
	if err != nil {
		log.Fatalf("Error getting sqlDB: %v", err)
	}
	defer sqlDB.Close()

	server.Use(middleware.DatabaseMiddleware(gormdb))
	server.Use(middleware.LoggerMiddleware())
	server.Run(os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"))
}
