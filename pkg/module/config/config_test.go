package config

import (
	"strings"
	"testing"
)

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
