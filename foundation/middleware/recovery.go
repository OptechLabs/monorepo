package middleware

import (
	"errors"
	"net/http"

	"github.com/OptechLabs/monorepo/foundation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		foundation.LoggerFrom(c).Error("panic occurred and recovered", zap.Any("error", err))
		pubErr := errors.New("An internal error has occurred. Contact Tech for more information")
		foundation.AbortWithError(c, http.StatusInternalServerError, pubErr)
	})
}
