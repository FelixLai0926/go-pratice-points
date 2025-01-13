package database

import (
	"fmt"
	"log"
	"points/pkg/module/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Instance *gorm.DB

func init() {
	envFilePath := config.GetEnvPath()

	if err := config.InitEnv(envFilePath); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	cfg, err := config.ParseEnv[PostgresConfig]()
	if err != nil {
		panic("Error transforming .env file to struct: " + err.Error())
	}

	dsn := GeneratePostgresDSN(cfg)

	Instance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := Instance.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Database connection pool configured successfully!")
}

func GeneratePostgresDSN(cfg *PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
}
