package main

import (
	"os"
	"points/pkg/module/config"
	"points/pkg/module/database"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	cfg := initConfig()

	gormdb := initDB(cfg)

	g := gen.NewGenerator(gen.Config{
		OutPath: os.Getenv("GEN_MODEL_OUT_PATH"),
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(gormdb)

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}

func initConfig() *database.PostgresConfig {
	envFilePath := config.GetEnvPath()
	if err := config.InitEnv(envFilePath); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	cfg, err := config.ParseEnv[database.PostgresConfig]()
	if err != nil {
		panic("Error transforming .env file to struct: " + err.Error())
	}
	return cfg
}

func initDB(cfg *database.PostgresConfig) *gorm.DB {
	dsn := database.GeneratePostgresDSN(cfg)

	gormdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	return gormdb
}
