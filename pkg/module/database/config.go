package database

type PostgresConfig struct {
	User     string `env:"POSTGRES_USER" required:"true"`
	Password string `env:"POSTGRES_PASSWORD" required:"true"`
	Host     string `env:"POSTGRES_HOST" required:"true"`
	Port     string `env:"POSTGRES_PORT" required:"true"`
	Database string `env:"POSTGRES_DB" required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE" required:"false"`
}
