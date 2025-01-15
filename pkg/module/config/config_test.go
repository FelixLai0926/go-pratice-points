package config

import (
	"os"
	"strings"
	"testing"
)

func TestGetRootPath(t *testing.T) {
	type args struct {
		environment []string
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "GetRootPath",
			wantErr: false,
		},
		{
			name:    "GetRootPath (error)",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				tmpDir, err := os.MkdirTemp("", "test")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(tmpDir)

				originalWd, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				defer os.Chdir(originalWd)

				if err := os.Chdir(tmpDir); err != nil {
					t.Fatal(err)
				}
				_, err = getRootPath()
				if err == nil {
					t.Error("getRootPath() expected error, got nil")
				}
				return
			}
			rootPath, err := getRootPath()
			if err == nil && tt.wantErr {
				t.Errorf("getRootPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rootPath == "" {
				t.Errorf("getRootPath() = %v, want not empty", rootPath)
			}
		})
	}
}

func TestGetEnvPath(t *testing.T) {
	type args struct {
		environment []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Production environment",
			args: args{environment: []string{"production"}},
			want: ".env.production",
		},
		{
			name: "Staging environment",
			args: args{environment: []string{"staging"}},
			want: ".env.staging",
		},
		{
			name: "Development environment (default)",
			args: args{environment: []string{}},
			want: ".env.example",
		},
		{
			name: "Multiple environments (use first)",
			args: args{environment: []string{"production", "staging"}},
			want: ".env.production",
		},
		{
			name: "Unknown environment",
			args: args{environment: []string{"unknown"}},
			want: ".env.unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEnvPath(tt.args.environment...)
			spiltResult := strings.Split(result, "/")
			fileName := spiltResult[len(spiltResult)-1]
			if fileName != tt.want {
				t.Errorf("GetEnvPath() = %v, want %v", fileName, tt.want)
			}
		})
	}
}

func TestInitEnv(t *testing.T) {
	type args struct {
		envFilePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test InitEnv",
			args:    args{envFilePath: "../../../configs/.env.example"},
			wantErr: false,
		},
		{
			name:    "Test InitEnv (unknown environment)",
			args:    args{envFilePath: "../../../configs/.env.unknown"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitEnv(tt.args.envFilePath); (err != nil) != tt.wantErr {
				t.Errorf("InitEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
