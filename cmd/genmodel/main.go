package main

import (
	"flag"
	"fmt"
	"log"

	"points/internal/domain/port"
	"points/internal/infrastructure"
	"points/internal/infrastructure/dbconnection"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	env := flag.String("env", "example", "specify the environment to use (example, development, production, etc.)")
	flag.Parse()
	viper, err := infrastructure.NewViperImpl(*env)
	if err != nil {
		log.Fatalf("Error initializing viper: %v", err)
	}

	defaults := infrastructure.NewDefaultSetterImpl()
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(viper, defaults, copier)

	dbConnection := dbconnection.NewPostgresConnection(config)

	if err != nil {
		log.Fatalf("Error initializing viper: %v", err)
	}

	db, err := dbConnection.InitPostgresDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer dbConnection.MustClose(db)

	log.Println("Generating models...")
	if err := generateModels(viper, db); err != nil {
		log.Fatalf("Error generating models: %v", err)
	}
	log.Println("Model generation completed.")
}

func generateModels(v port.SettingsManager, gormdb *gorm.DB) error {
	daoPath := v.GetString("gen.GEN_DAO_PATH")
	if daoPath == "" {
		return fmt.Errorf("GEN_DAO_PATH is not set")
	}

	modelPath := v.GetString("gen.GEN_MODEL_OUT_PATH")
	if modelPath == "" {
		return fmt.Errorf("GEN_MODEL_OUT_PATH is not set")
	}

	generator := gen.NewGenerator(gen.Config{
		OutPath:      daoPath,
		ModelPkgPath: modelPath,
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
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
