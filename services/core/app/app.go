package app

import (
	"context"
	"net/http"

	"github.com/OptechLabs/monorepo/foundation"
	config "github.com/OptechLabs/monorepo/helpers/config"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/source/google_cloud_storage"
)

func New(
	ctx context.Context,
	logger foundation.Logger,
	config config.Config,
) (app *foundation.Foundation, shutdown func() error, err error) {

	app = foundation.New(foundation.Options{
		Environment:     config.Environment,
		HTTPPort:        config.HTTPServerConfig.Port,
		StartHTTPServer: config.HTTPServerConfig.Port != "",
		Logger:          logger,
	})

	noAuthRoutes := app.HTTPRouter.Group("/")
	{
		noAuthRoutes.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	return app, func() error {
		// waiting for shutdown signal
		<-ctx.Done()
		// shutdown received, stopping...
		return nil
	}, nil
}
