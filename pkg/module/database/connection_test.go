package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePostgresDSN(t *testing.T) {
	dsn, err := GeneratePostgresDSN(nil)
	assert.Error(t, err, "nil config should return error")
	assert.Equal(t, "", dsn)

	cfgMissing := &PostgresConfig{
		User:     "",
		Password: "",
		Host:     "",
		Port:     "",
		Database: "",
	}
	dsn, err = GeneratePostgresDSN(cfgMissing)
	assert.Error(t, err, "missing required fields should return error")
	assert.Equal(t, "", dsn)

	cfgValid := &PostgresConfig{
		User:     "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     "5432",
		Database: "points",
		SSLMode:  "",
	}
	dsn, err = GeneratePostgresDSN(cfgValid)
	assert.NoError(t, err)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/points?sslmode=disable", dsn)

	cfgValid.SSLMode = "enable"
	dsn, err = GeneratePostgresDSN(cfgValid)
	assert.NoError(t, err)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/points?sslmode=enable", dsn)
}

func TestInitDatabase_InvalidConfig(t *testing.T) {
	db, err := InitDatabase(nil)
	assert.Error(t, err, "nil config should return error")
	assert.Nil(t, db)

	cfgMissing := &PostgresConfig{
		User:     "postgres",
		Password: "postgres",
		Host:     "",
		Port:     "5432",
		Database: "points",
		SSLMode:  "",
	}
	db, err = InitDatabase(cfgMissing)
	assert.Error(t, err, "missing required fields should return error")
	assert.Nil(t, db)
}

func TestClose_NilDB(t *testing.T) {
	assert.Panics(t, func() { _ = Close(nil) }, "Close(nil) should panic")
}
