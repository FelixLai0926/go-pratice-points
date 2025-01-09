package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	envFile := ".env.example"
	if len(os.Args) > 1 {
		envFile = os.Args[1]
	} else if envEnv := os.Getenv("ENV_FILE"); envEnv != "" {
		envFile = envEnv
	}

	if err := godotenv.Load(envFile); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	fmt.Print(dsn)
	gormdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./pkg/model",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(gormdb)

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
