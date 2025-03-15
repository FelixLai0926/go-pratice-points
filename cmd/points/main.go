package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"points/internal/di"

	"go.uber.org/fx"
)

func main() {
	env := flag.String("env", "example", "specify the environment to use (example, development, production, etc.)")
	flag.Parse()

	fmt.Println("Using environment:", *env)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		fx.Invoke(di.StartServer),
	)

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-app.Done()
	fmt.Println("Shutting down gracefully...")
}
