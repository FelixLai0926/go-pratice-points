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
			envFilePath := config.GetEnvPath(tt.environment)
			print(envFilePath)
			if err := config.InitEnv(envFilePath); err != nil {
				t.Fatalf("error loading .env file: %v", err)
			}

			cfg, err := config.ParseEnv[PostgresConfig]()
			if err != nil {
				t.Fatalf("error transforming .env file to struct: %v", err)
			}
			if err := InitDatabase(cfg); (err != nil) != tt.wantErr {
				t.Fatalf("InitDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePostgresDSN(tt.args.cfg); got != tt.want {
				t.Errorf("GeneratePostgresDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
