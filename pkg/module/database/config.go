package database

type PostgresConfig struct {
	User        string `env:"POSTGRES_USER" required:"true"`
	Password    string `env:"POSTGRES_PASSWORD" required:"true"`
	Host        string `env:"POSTGRES_HOST" required:"true"`
	Port        string `env:"POSTGRES_PORT" required:"true"`
	Database    string `env:"POSTGRES_DB" required:"true"`
	SSLMode     string `env:"POSTGRES_SSLMODE" required:"false"`
	MaxOpenConn int    `env:"POSTGRES_MAX_OPEN_CONNS" required:"false" default:"100"`
	MaxIdleConn int    `env:"POSTGRES_MAX_IDLE_CONNS" required:"false" default:"10"`
	ConnMaxLife int    `env:"POSTGRES_CONN_MAX_LIFE" required:"false" default:"30"`
}
