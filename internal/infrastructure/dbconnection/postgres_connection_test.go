package dbconnection

import (
	"context"
	"errors"
	"testing"
	"time"

	"points/test/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGetPostgresConfig_UnmarshalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mock.NewMockConfig(ctrl)
	dummySettingsManager := mock.NewMockSettingsManager(ctrl)
	mockCfg.EXPECT().Sub("postgres").Return(dummySettingsManager).Times(1)
	dummySettingsManager.EXPECT().Unmarshal(gomock.Any()).Return(errors.New("unmarshal error"))

	conn := NewPostgresConnection(mockCfg)
	_, err := conn.getPostgresConfig()
	if err == nil || err.Error() != "unmarshal error" {
		t.Errorf("expected unmarshal error, got: %v", err)
	}
}

func TestGetPostgresConfig_SetDefaultError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mock.NewMockConfig(ctrl)
	validCfg := PostgresConfig{
		User:     "user",
		Password: "pass",
		Host:     "localhost",
		Port:     "5432",
		Database: "db",
		SSLMode:  "disable",
	}

	dummySettingsManager := mock.NewMockSettingsManager(ctrl)
	mockCfg.EXPECT().Sub("postgres").Return(dummySettingsManager).Times(1)
	dummySettingsManager.EXPECT().Unmarshal(gomock.Any()).DoAndReturn(func(out interface{}) error {
		*(out.(*PostgresConfig)) = validCfg
		return nil
	})
	mockCfg.EXPECT().SetDefault(gomock.Any()).Return(errors.New("setdefault error"))

	conn := NewPostgresConnection(mockCfg)
	_, err := conn.getPostgresConfig()
	if err == nil || err.Error() != "setdefault error" {
		t.Errorf("expected setdefault error, got: %v", err)
	}
}

func TestInitPostgresDatabase_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mock.NewMockConfig(ctrl)
	invalidCfg := PostgresConfig{}
	dummySettingsManager := mock.NewMockSettingsManager(ctrl)
	mockCfg.EXPECT().Sub("postgres").Return(dummySettingsManager).Times(1)
	dummySettingsManager.EXPECT().Unmarshal(gomock.Any()).DoAndReturn(func(out interface{}) error {
		*(out.(*PostgresConfig)) = invalidCfg
		return nil
	})
	mockCfg.EXPECT().SetDefault(gomock.Any()).Return(nil)

	conn := NewPostgresConnection(mockCfg)
	_, err := conn.InitPostgresDatabase()
	if err == nil {
		t.Errorf("expected validation error due to missing required fields, got nil")
	}
}

func TestInitPostgresDatabase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Helper()
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err, "failed to start postgres container")

	port, err := container.MappedPort(ctx, "5432")

	mockCfg := mock.NewMockConfig(ctrl)
	validCfg := PostgresConfig{
		User:        "test",
		Password:    "test",
		Host:        "localhost",
		Port:        port.Port(),
		Database:    "testdb",
		SSLMode:     "disable",
		MaxOpenConn: 10,
		MaxIdleConn: 5,
		ConnMaxLife: 30,
	}

	dummySettingsManager := mock.NewMockSettingsManager(ctrl)
	mockCfg.EXPECT().Sub("postgres").Return(dummySettingsManager).Times(1)
	dummySettingsManager.EXPECT().Unmarshal(gomock.Any()).DoAndReturn(func(out interface{}) error {
		*(out.(*PostgresConfig)) = validCfg
		return nil
	})
	mockCfg.EXPECT().SetDefault(gomock.Any()).Return(nil)

	conn := NewPostgresConnection(mockCfg)
	db, err := conn.InitPostgresDatabase()
	if err != nil {
		t.Fatalf("expected successful init, got error: %v", err)
	}

	if err := conn.Close(db); err != nil {
		t.Errorf("failed to close database: %v", err)
	}
}

func TestClose_NilDB(t *testing.T) {
	conn := NewPostgresConnection(nil)
	err := conn.Close(nil)
	if err == nil {
		t.Errorf("expected error when closing nil db, got nil")
	}
}
