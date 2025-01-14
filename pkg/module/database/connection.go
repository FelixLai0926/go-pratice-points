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

func InitDatabase() error {
	envFilePath := config.GetEnvPath()

	if err := config.InitEnv(envFilePath); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	cfg, err := config.ParseEnv[PostgresConfig]()
	if err != nil {
		return fmt.Errorf("error transforming .env file to struct: %v", err)
	}

	dsn := GeneratePostgresDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLife) * time.Minute)

	log.Println("Database connection pool configured successfully!")

	Instance = db
	return nil
}

func GeneratePostgresDSN(cfg *PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
}
