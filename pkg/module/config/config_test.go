package config

import (
	"os"
	"strings"
	"testing"
)

func TestGetRootPath(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Valid root path", false},
		{"Missing go.mod", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				tmpDir, _ := os.MkdirTemp("", "test")
				defer os.RemoveAll(tmpDir)

				originalWd, _ := os.Getwd()
				defer os.Chdir(originalWd)

				os.Chdir(tmpDir)
			}

			rootPath, err := GetRootPath()
			if (err == nil && tt.wantErr) || (err != nil && !tt.wantErr) {
				t.Fatalf("GetRootPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rootPath == "" && !tt.wantErr {
				t.Fatalf("GetRootPath() returned empty path")
			}
		})
	}
}

func TestGetEnvPath(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"Valid production environment", []string{"production"}, ".env.production"},
		{"Valid staging environment", []string{"staging"}, ".env.staging"},
		{"Default environment", []string{}, ".env.example"},
		{"Fallback to first environment", []string{"production", "staging"}, ".env.production"},
		{"Unknown environment", []string{"unknown"}, ".env.unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := GetEnvPath(tt.args...)
			if !strings.HasSuffix(result, tt.want) {
				t.Errorf("GetEnvPath() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestInitEnv(t *testing.T) {
	tests := []struct {
		name    string
		envFile string
		wantErr bool
	}{
		{"Valid .env file", "POSTGRES_USER=postgres\nPOSTGRES_PASSWORD=secret\n", false},
		{"Missing .env file", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tmpFile *os.File
			var err error

			if tt.envFile != "" {
				tmpFile, err = os.CreateTemp("", "env.unittest")
				if err != nil {
					t.Fatalf("faild to create temp file: %v", err)
				}
				defer os.Remove(tmpFile.Name())

				if _, err := tmpFile.WriteString(tt.envFile); err != nil {
					t.Fatalf("faild to write to temp file: %v", err)
				}
				tmpFile.Close()
			}

			envFilePath := ""
			if tmpFile != nil {
				envFilePath = tmpFile.Name()
			}

			err = InitEnv(envFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
