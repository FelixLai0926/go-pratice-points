package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"points/pkg/module/config"
	"points/pkg/module/database"
	"strings"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	env := flag.String("env", "example", "specify the environment to use (example, development, production, etc.)")
	flag.Parse()

	cfg, err := loadDatabaseConfig(*env)
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

func loadDatabaseConfig(environment string) (*database.PostgresConfig, error) {
	envData, err := config.GetEnvData(environment)
	if err != nil {
		return nil, fmt.Errorf("Error getting env data: %v", err)
	}
	if err := config.InitEmbeddedEnv(envData); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	cfg, err := config.ParseEnv[database.PostgresConfig]()
	if err != nil {
		return nil, fmt.Errorf("Error transforming .env file to struct: %v", err)
	}
	return cfg, nil
}

func generateModels(gormdb *gorm.DB) error {
	daoPath := os.Getenv("GEN_DAO_PATH")
	if daoPath == "" {
		return fmt.Errorf("GEN_DAO_PATH is not set")
	}

	modelPath := os.Getenv("GEN_MODEL_OUT_PATH")
	if modelPath == "" {
		return fmt.Errorf("GEN_MODEL_OUT_PATH is not set")
	}

	generator := gen.NewGenerator(gen.Config{
		OutPath:      daoPath,
		ModelPkgPath: modelPath,
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	generator.WithModelNameStrategy(func(tableName string) string {
		if tableName == "transaction" {
			return "TransactionDAO"
		}
		return strings.Title(tableName)
	})

	generator.WithImportPkgPath("github.com/shopspring/decimal")
	generator.WithDataTypeMap(map[string]func(detail gorm.ColumnType) (dataType string){
		"decimal": func(columnType gorm.ColumnType) string {
			return "decimal.Decimal"
		},
		"numeric": func(columnType gorm.ColumnType) string {
			return "decimal.Decimal"
		},
	})

	generator.UseDB(gormdb)

	generator.ApplyBasic(generator.GenerateAllTable()...)

	generator.Execute()

	return nil
}
