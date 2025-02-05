package main

import (
	"fmt"
	"log"
	"os"
	"points/pkg/module/config"
	"points/pkg/module/database"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	cfg, err := loadDatabaseConfig()
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	gormdb, err := database.InitDatabase(cfg)
	defer database.Close(gormdb)

	log.Println("Generating models...")
	if err := generateModels(gormdb); err != nil {
		log.Fatalf("Error generating models: %v", err)
	}
	log.Println("Model generation completed.")
}

func loadDatabaseConfig() (*database.PostgresConfig, error) {
	envFilePath, err := config.GetEnvPath()
	if err != nil {
		return nil, fmt.Errorf("Error getting .env file path: %v", err)
	}
	if err := config.InitEnv(envFilePath); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	cfg, err := config.ParseEnv[database.PostgresConfig]()
	if err != nil {
		return nil, fmt.Errorf("Error transforming .env file to struct: %v", err)
	}
	return cfg, nil
}

func generateModels(gormdb *gorm.DB) error {
	outPutPath := os.Getenv("GEN_MODEL_OUT_PATH")
	if outPutPath == "" {
		return fmt.Errorf("GEN_MODEL_OUT_PATH is not set")
	}

	generator := gen.NewGenerator(gen.Config{
		OutPath: outPutPath,
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	generator.UseDB(gormdb)

	generator.ApplyBasic(generator.GenerateAllTable()...)

	generator.Execute()

	return nil
}
