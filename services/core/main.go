package main

import (
	"log"
	"os"

	"github.com/OptechLabs/monorepo/foundation"
	config "github.com/OptechLabs/monorepo/helpers/config"
	"github.com/OptechLabs/monorepo/services/gateway/app"
	"go.uber.org/zap"

	"golang.org/x/sync/errgroup"
)

func main() {
	appConfig := loadConfigurations()
	if appConfig.Environment != "development" {
		appConfig.GRPCServerConfig.Port = os.Getenv("PORT")
	}

	logger, _ := foundation.NewDefaultLogger(appConfig.Environment)
	ctx, stop := foundation.ContextWithCancel()

	app, shutdown, err := app.New(ctx, logger, appConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	g := new(errgroup.Group)
	g.Go(func() error {
		return app.RunWithContext(ctx, stop)
	})
	g.Go(func() error {
		return shutdown()
	})
	if err := g.Wait(); err != nil {
		logger.Error("failed running or shutting down example app", zap.Error(err))
		return
	}
}

func loadConfigurations() (appConfig config.Config) {
	var err error
	cfgFile := os.Getenv("LOCAL_CONFIG_FILE")
	configStr := os.Getenv("config")

	if cfgFile == "" && configStr == "" {
		log.Fatal("Either the LOCAL_CONFIG_FILE or config environment variable must be set.")
	}

	if configStr != "" {
		appConfig, err = config.LoadFromString(configStr, true)
		if err != nil {
			log.Fatal("failed to load config", err)
		}
	}

	if os.Getenv("LOCAL_CONFIG_FILE") != "" {
		appConfig, err = config.LoadFromFile(cfgFile, true)
		if err != nil {
			log.Fatal("failed to load config", err)
		}
	}

	return appConfig
}
