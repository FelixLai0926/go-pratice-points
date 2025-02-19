package main

import (
	"log"
	"os"
	"points/pkg/middleware"
	"points/pkg/module/cache"
	"points/pkg/module/config"
	"points/pkg/module/database"
	"points/pkg/module/logger"
	"points/pkg/router"

	"github.com/gin-gonic/gin"
)

func main() {
	env := "example"
	if err := logger.InitLogger(env); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.SyncLogger()

	envFilePath, err := config.GetEnvPath(env)
	if err != nil {
		log.Fatalf("Error getting .env file path: %v", err)
	}
	if err := config.InitEnv(envFilePath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//Postgres
	postgresConfig, err := config.ParseEnv[database.PostgresConfig]()
	if err != nil {
		log.Fatalf("Error transforming .env file to struct: %v", err)
	}

	gormdb, err := database.InitDatabase(postgresConfig)
	if err != nil {
		log.Fatalf("Error initializing database:: %v", err)
	}

	sqlDB, err := gormdb.DB()
	if err != nil {
		log.Fatalf("Error getting sqlDB: %v", err)
	}
	defer sqlDB.Close()

	//redis
	redisConfig, err := config.ParseEnv[cache.RedisConfig]()
	log.Printf("redis config: %v", redisConfig)
	if err != nil {
		log.Fatalf("Error transforming .env file to struct: %v", err)
	}

	redisClient, err := cache.InitRedisClient(redisConfig)
	if err != nil {
		log.Fatalf("Error initializing redis:: %v", err)
	}
	defer redisClient.Close()

	server := gin.Default()
	server.Use(middleware.LoggerMiddleware())
	server.Use(middleware.ErrorHandlerMiddleware())
	router.RegisterTestRoutes(server)
	router.RegisterUserRoutes(server, gormdb, redisClient)
	server.Run(os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"))
}
