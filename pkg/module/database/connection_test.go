package database

import (
	"points/pkg/module/config"
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.NoError(t, err, "Unexpected error when getting .env file path")

			err = config.InitEnv(envFilePath)
			if tt.wantErr {
				assert.Error(t, err, "Expected error loading .env file for environment 'unknown'")
				return
			} else {
				assert.NoError(t, err, "Unexpected error when loading .env file")
			}

			cfg, err := config.ParseEnv[PostgresConfig]()
			assert.NoError(t, err, "Unexpected error when parsing .env file to struct")

			_, err = InitDatabase(cfg)
			if tt.wantErr {
				assert.Error(t, err, "Expected error when initializing database")
			} else {
				assert.NoError(t, err, "Unexpected error when initializing database")
			}
		})
	}
}

func TestGeneratePostgresDSN(t *testing.T) {
	type args struct {
		cfg *PostgresConfig
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
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
			want:    "postgres://postgres:postgres@localhost:5432/points?sslmode=disable",
			wantErr: false,
		},
		{
			name: "Valid DSN with SSLMode enable",
			args: args{
				cfg: &PostgresConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "postgres",
					Database: "points",
					SSLMode:  "enable",
				},
			},
			want:    "postgres://postgres:postgres@localhost:5432/points?sslmode=enable",
			wantErr: false,
		},
		{
			name: "Nil config",
			args: args{
				cfg: nil,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePostgresDSN(tt.args.cfg)
			if tt.wantErr {
				assert.Error(t, err, "Expected error generating DSN")
			} else {
				assert.NoError(t, err, "Unexpected error generating DSN")
				assert.Equal(t, tt.want, got, "DSN mismatch")
			}
		})
	}
}
