package config

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"points/pkg/module/database"

	"github.com/stretchr/testify/assert"
)

func Test_setFieldValue(t *testing.T) {
	type TestStruct struct {
		IntField    int
		FloatField  float64
		BoolField   bool
		StringField string
	}
	type args struct {
		field reflect.Value
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test setFieldValue (string)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("StringField"),
				value: "test",
			},
			wantErr: false,
		},
		{
			name: "Test setFieldValue (int)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("IntField"),
				value: "10",
			},
			wantErr: false,
		},
		{
			name: "Test setFieldValue (float)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("FloatField"),
				value: "10.5",
			},
			wantErr: false,
		},
		{
			name: "Test setFieldValue (bool)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("BoolField"),
				value: "true",
			},
			wantErr: false,
		},
		{
			name: "Test setFieldValue (invalid int)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("IntField"),
				value: "test",
			},
			wantErr: true,
		},
		{
			name: "Test setFieldValue (invalid float)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("FloatField"),
				value: "test",
			},
			wantErr: true,
		},
		{
			name: "Test setFieldValue (invalid bool)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("BoolField"),
				value: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setFieldValue(tt.args.field, tt.args.value)
			if tt.wantErr {
				assert.Error(t, err, "Expected error for value: %s", tt.args.value)
			} else {
				assert.NoError(t, err, "Did not expect error for value: %s", tt.args.value)
			}
		})
	}
}

func TestParseEnv(t *testing.T) {
	tests := []struct {
		name    string
		want    *database.PostgresConfig
		wantErr bool
		osEnv   map[string]string
	}{
		{
			name: "Complete environment variables",
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
			name: "Missing optional variables (use default values)",
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
			name:    "Empty environment variables (return error)",
			want:    nil,
			wantErr: true,
			osEnv:   map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.osEnv {
				os.Setenv(key, value)
				t.Cleanup(func() {
					os.Unsetenv(key)
				})
			}
			got, err := ParseEnv[database.PostgresConfig]()
			if tt.wantErr {
				assert.Error(t, err, "Expected error when environment variables are missing")
			} else {
				assert.NoError(t, err, "Unexpected error when parsing environment variables")
				assert.Equal(t, tt.want, got, "Parsed config does not match expected config")
			}
		})
	}
}

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
				tmpDir, err := os.MkdirTemp("", "test")
				assert.NoError(t, err, "Failed to create temp directory")
				originalWd, err := os.Getwd()
				assert.NoError(t, err, "Failed to get current working directory")
				err = os.Chdir(tmpDir)
				assert.NoError(t, err, "Failed to change directory to temp directory")
				t.Cleanup(func() {
					os.Chdir(originalWd)
					os.RemoveAll(tmpDir)
				})
			}

			rootPath, err := GetRootPath()
			if tt.wantErr {
				assert.Error(t, err, "Expected error when go.mod is missing")
			} else {
				assert.NoError(t, err, "Unexpected error when getting root path")
				assert.NotEmpty(t, rootPath, "Root path should not be empty")
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
			result, err := GetEnvPath(tt.args...)
			assert.NoError(t, err, "Unexpected error in GetEnvPath")
			assert.True(t, strings.HasSuffix(result, tt.want), "GetEnvPath() = %v, want suffix %v", result, tt.want)
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
			var tmpFilePath string
			if tt.envFile != "" {
				tmpFile, err := os.CreateTemp("", "env.unittest")
				assert.NoError(t, err, "Failed to create temp file")
				_, err = tmpFile.WriteString(tt.envFile)
				assert.NoError(t, err, "Failed to write to temp file")
				tmpFilePath = tmpFile.Name()
				tmpFile.Close()
				t.Cleanup(func() {
					os.Remove(tmpFilePath)
				})
			}

			err := InitEnv(tmpFilePath)
			if tt.wantErr {
				assert.Error(t, err, "Expected error when loading missing .env file")
			} else {
				assert.NoError(t, err, "Unexpected error when loading valid .env file")
			}
		})
	}
}
