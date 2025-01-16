package database

import (
	"points/pkg/module/config"
	"testing"
)

func TestInitDatabase(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		wantErr     bool
	}{
		{
			name:        "TestInitDatabase",
			environment: "example",
			wantErr:     false,
		},
		{
			name:        "TestInitDatabase (error)",
			environment: "unknown",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envFilePath, err := config.GetEnvPath(tt.environment)
			checkError(t, err, tt.wantErr, "error getting .env file path", false)
			if tt.wantErr && err != nil {
				return
			}

			err = config.InitEnv(envFilePath)
			checkError(t, err, tt.wantErr, "error loading .env file", false)
			if tt.wantErr && err != nil {
				return
			}

			cfg, err := config.ParseEnv[PostgresConfig]()
			checkError(t, err, tt.wantErr, "error transforming .env file to struct", false)
			if tt.wantErr && err != nil {
				return
			}

			err = InitDatabase(cfg)
			checkError(t, err, tt.wantErr, "error initializing database", true)
		})
	}
}

func TestGeneratePostgresDSN(t *testing.T) {
	type args struct {
		cfg *PostgresConfig
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid DSN with SSLMode empty",
			args: args{
				cfg: &PostgresConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "postgres",
					Database: "points",
				},
			},
			want: "postgres://postgres:postgres@localhost:5432/points?sslmode=",
		},
		{
			name: "Valid DSN with SSLMode disable",
			args: args{
				cfg: &PostgresConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "postgres",
					Database: "points",
					SSLMode:  "disable",
				},
			},
			want: "postgres://postgres:postgres@localhost:5432/points?sslmode=disable",
		},
		{
			name: "Nil config",
			args: args{
				cfg: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePostgresDSN(tt.args.cfg); got != tt.want {
				t.Errorf("GeneratePostgresDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func checkError(t *testing.T, err error, wantErr bool, errMsg string, final bool) {
	if err != nil && !wantErr {
		t.Fatalf("%s: %v", errMsg, err)
	}
	if err == nil && wantErr && final {
		t.Fatalf("%s: expected error but got nil", errMsg)
	}
}
