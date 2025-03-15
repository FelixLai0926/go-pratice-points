package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"points/internal/di"
	"points/internal/domain/port"

	"go.uber.org/fx"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	env := flag.String("env", "example", "specify the environment to use (example, development, production, etc.)")
	flag.Parse()

	app := fx.New(
		fx.Supply(*env),
		di.SettingManagerModule,
		di.DefaultsModule,
		di.CopierModule,
		di.ConfigModule,
		di.LoggerModule,
		di.DatabaseModule,
		di.ApplicationModule,
		di.HTTPModule,
		fx.Invoke(runGenerateModels),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := app.Stop(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func runGenerateModels(config port.Config, db *gorm.DB) {
	log.Println("Generating models...")
	if err := generateModels(config, db); err != nil {
		log.Fatalf("Error generating models: %v", err)
	}
	log.Println("Model generation completed.")
}

func generateModels(config port.Config, gormdb *gorm.DB) error {
	daoPath := config.GetString("gen.GEN_DAO_PATH")
	if daoPath == "" {
		return fmt.Errorf("GEN_DAO_PATH is not set")
	}

	modelPath := config.GetString("gen.GEN_MODEL_OUT_PATH")
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
