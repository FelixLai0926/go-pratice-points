package dbconnection

import (
	"fmt"
	"log"
	"points/internal/domain/port"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresConfig struct {
	User        string `mapstructure:"POSTGRES_USER" validate:"required"`
	Password    string `mapstructure:"POSTGRES_PASSWORD" validate:"required"`
	Host        string `mapstructure:"POSTGRES_HOST" validate:"required"`
	Port        string `mapstructure:"POSTGRES_PORT" validate:"required"`
	Database    string `mapstructure:"POSTGRES_DB" validate:"required"`
	SSLMode     string `mapstructure:"POSTGRES_SSLMODE" validate:"omitempty"`
	MaxOpenConn int    `mapstructure:"POSTGRES_MAX_OPEN_CONNS" validate:"omitempty"`
	MaxIdleConn int    `mapstructure:"POSTGRES_MAX_IDLE_CONNS" validate:"omitempty"`
	ConnMaxLife int    `mapstructure:"POSTGRES_CONN_MAX_LIFE" validate:"omitempty"`
}

type PostgresConnection struct {
	config   port.Config
	validate *validator.Validate
}

func NewPostgresConnection(config port.Config) *PostgresConnection {
	return &PostgresConnection{config: config, validate: validator.New()}
}

func (d *PostgresConnection) InitPostgresDatabase() (*gorm.DB, error) {
	cfg, err := d.getPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres config: %w", err)
	}

	err = d.validate.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to validate postgres config: %w", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	connMaxLife := cfg.ConnMaxLife
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLife) * time.Minute)

	return db, nil
}

func (d *PostgresConnection) Close(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("failed to close database, got nil database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

func (d *PostgresConnection) MustClose(db *gorm.DB) {
	if err := d.Close(db); err != nil {
		log.Printf("failed to close database: %v", err)
	}
}

func (d *PostgresConnection) getPostgresConfig() (*PostgresConfig, error) {
	var cfg PostgresConfig
	v := d.config.Sub("postgres")
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := d.config.SetDefault(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
