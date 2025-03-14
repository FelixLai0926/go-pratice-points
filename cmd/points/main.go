package main

import (
	"flag"
	"log"
	"os"

	"points/internal/adapter/http/middleware"
	"points/internal/adapter/http/router"
	"points/internal/infrastructure"
	"points/internal/infrastructure/dbconnection"

	"github.com/gin-gonic/gin"
)

func main() {
	env := flag.String("env", "example", "specify the environment to use (example, development, production, etc.)")
	flag.Parse()

	logger, err := infrastructure.NewZapLogger(*env)
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logger.Sync()

	viper, err := infrastructure.NewViperImpl(*env)
	if err != nil {
		log.Fatalf("Error initializing viper: %v", err)
	}
	defaults := infrastructure.NewDefaultSetterImpl()
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(viper, defaults, copier)

	//Postgres
	postgresConnection := dbconnection.NewPostgresConnection(config)
	gormdb, err := postgresConnection.InitPostgresDatabase()
	if err != nil {
		log.Fatalf("Error initializing database:: %v", err)
	}

	defer postgresConnection.Close(gormdb)

	//redis
	redisConnection := dbconnection.NewRedisConnection(config)
	redisClient, err := redisConnection.InitRedisDatabase()
	if err != nil {
		log.Fatalf("Error initializing redis:: %v", err)
	}

	server := gin.Default()
	server.Use(middleware.LoggerMiddleware(logger))
	server.Use(middleware.ErrorHandlerMiddleware(logger))
	router.RegisterTestRoutes(server)
	router.RegisterUserRoutes(server, gormdb, redisClient, config)
	server.Run(os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"))
}
