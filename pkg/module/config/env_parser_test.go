package config

import (
	"fmt"
	"os"
	"points/pkg/module/database"
	"reflect"
	"testing"
)

func TestParseEnv(t *testing.T) {
	tests := []struct {
		name    string
		want    *database.PostgresConfig
		wantErr bool
		osEnv   map[string]string
	}{
		{
			name: "Test ParseEnv",
			want: &database.PostgresConfig{
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
			want: &database.PostgresConfig{
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
			got, err := ParseEnv[database.PostgresConfig]()
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
