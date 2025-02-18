package test

import (
	"context"
	"fmt"
	"points/pkg/models/orm"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/docker/go-connections/nat"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDummyDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&orm.Account{})
	assert.NoError(t, err, "failed to migrate database schema")
	err = db.AutoMigrate(&orm.Transaction{})
	assert.NoError(t, err, "failed to migrate database schema")
	err = db.AutoMigrate(&orm.TransactionEvent{})
	assert.NoError(t, err, "failed to migrate database schema")
	assert.NoError(t, err)
	return db
}

func NewDummyRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return mr, client
}

func SetupRedisContainer(t *testing.T) (string, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:7.0",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort(nat.Port("6379/tcp")).WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	port, err := container.MappedPort(ctx, "6379")
	assert.NoError(t, err)
	dsn := fmt.Sprintf("redis://localhost:%s/0", port.Port())
	cleanup := func() {
		container.Terminate(ctx)
	}
	return dsn, cleanup
}

func SetupAccounts(t *testing.T, db *gorm.DB) {
	fromAcc := orm.Account{
		UserID:           1,
		AvailableBalance: 1000.0,
		ReservedBalance:  0.0,
	}
	toAcc := orm.Account{
		UserID:           2,
		AvailableBalance: 500.0,
		ReservedBalance:  0.0,
	}
	if err := db.Create(&fromAcc).Error; err != nil {
		t.Fatalf("failed to create from account: %v", err)
	}
	if err := db.Create(&toAcc).Error; err != nil {
		t.Fatalf("failed to create to account: %v", err)
	}
}

func NewTestContainerDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn, cleanup := setupPostgresContainer(t)
	t.Cleanup(cleanup)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err, "failed to open gorm db")

	sqlDB, err := db.DB()
	assert.NoError(t, err)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)

	err = db.AutoMigrate(&orm.Account{})
	assert.NoError(t, err, "failed to migrate database schema")
	err = db.AutoMigrate(&orm.Transaction{})
	assert.NoError(t, err, "failed to migrate database schema")
	err = db.AutoMigrate(&orm.TransactionEvent{})
	assert.NoError(t, err, "failed to migrate database schema")

	return db
}

func setupPostgresContainer(t *testing.T) (string, func()) {
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
	assert.NoError(t, err, "failed to map container port")

	dsn := fmt.Sprintf("host=localhost user=test password=test dbname=testdb port=%s sslmode=disable", port.Port())

	cleanup := func() {
		err := container.Terminate(ctx)
		assert.NoError(t, err, "failed to terminate container")
	}

	return dsn, cleanup
}
