package config

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestParseEnv(t *testing.T) {
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

	tests := []struct {
		name    string
		want    *PostgresConfig
		wantErr bool
		osEnv   map[string]string
	}{
		{
			name: "Test ParseEnv",
			want: &PostgresConfig{
				User:        "postgres",
				Password:    "pgadmin1234",
				Host:        "127.0.0.1",
				Port:        "5432",
				Database:    "points",
				SSLMode:     "disable",
				MaxOpenConn: 200,
				MaxIdleConn: 20,
				ConnMaxLife: 60,
			},
			wantErr: false,
			osEnv: map[string]string{
				"POSTGRES_USER":           "postgres",
				"POSTGRES_PASSWORD":       "pgadmin1234",
				"POSTGRES_HOST":           "127.0.0.1",
				"POSTGRES_PORT":           "5432",
				"POSTGRES_DB":             "points",
				"POSTGRES_SSLMODE":        "disable",
				"POSTGRES_MAX_OPEN_CONNS": "200",
				"POSTGRES_MAX_IDLE_CONNS": "20",
				"POSTGRES_CONN_MAX_LIFE":  "60",
			},
		},
		{
			name: "Test ParseEnv (default values)",
			want: &PostgresConfig{
				User:        "postgres",
				Password:    "pgadmin1234",
				Host:        "127.0.0.1",
				Port:        "5432",
				Database:    "points",
				SSLMode:     "disable",
				MaxOpenConn: 100,
				MaxIdleConn: 10,
				ConnMaxLife: 30,
			},
			wantErr: false,
			osEnv: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "pgadmin1234",
				"POSTGRES_HOST":     "127.0.0.1",
				"POSTGRES_PORT":     "5432",
				"POSTGRES_DB":       "points",
				"POSTGRES_SSLMODE":  "disable",
			},
		},
		{
			name:    "Test ParseEnv (empty environment)",
			want:    nil,
			wantErr: true,
			osEnv:   map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.osEnv {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}
			got, err := ParseEnv[PostgresConfig]()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
