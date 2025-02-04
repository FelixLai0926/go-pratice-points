package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(cfg *PostgresConfig) (*gorm.DB, error) {
	dsn, err := GeneratePostgresDSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DSN: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %v", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLife) * time.Minute)

	return db, nil
}

func Close(gormdb *gorm.DB) error {
	sqlDB, err := gormdb.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}
	return sqlDB.Close()
}

func GeneratePostgresDSN(cfg *PostgresConfig) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("GeneratePostgresDSN expects a non-nil PostgresConfig")
	}

	if cfg.User == "" || cfg.Password == "" || cfg.Host == "" || cfg.Port == "" || cfg.Database == "" {
		missingFields := []string{}
		if cfg.User == "" {
			missingFields = append(missingFields, "User")
		}
		if cfg.Password == "" {
			missingFields = append(missingFields, "Password")
		}
		if cfg.Host == "" {
			missingFields = append(missingFields, "Host")
		}
		if cfg.Port == "" {
			missingFields = append(missingFields, "Port")
		}
		if cfg.Database == "" {
			missingFields = append(missingFields, "Database")
		}
		return "", fmt.Errorf("missing required fields in PostgresConfig: %v", missingFields)
	}

	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
	return dsn, nil
}
