package middleware

import (
	"time"

	"github.com/OptechLabs/monorepo/foundation"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// RequestIDField is the header field name of the request ID
	RequestIDField string

	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

const RequestIDField = "X-Request-ID"

// Logger returns a middleware that that will write structured http request logs with zap
func Logger(logger foundation.Logger) gin.HandlerFunc {
	return LoggerWithConfig(logger, LoggerConfig{
		RequestIDField: RequestIDField,
	})
}

func LoggerWithConfig(logger foundation.Logger, conf LoggerConfig) gin.HandlerFunc {
	notlogged := conf.SkipPaths

	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		reqID := c.Request.Header.Get(conf.RequestIDField)
		if reqID == "" {
			reqID = uuid.Must(uuid.NewV4()).String()
		}
		c.Header(conf.RequestIDField, reqID)
		c.Set(foundation.RequestIDKey, reqID)
		c.Set(foundation.LoggerKey, logger.With(zap.String("request_id", reqID)))

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			latency := time.Since(start)
			if latency > time.Minute {
				latency = latency.Truncate(time.Second)
			}

			if raw != "" {
				path = path + "?" + raw
			}

			contentType := c.ContentType()
			if contentType == "" {
				contentType = c.Writer.Header().Get("Content-Type")
			}

			msg := "[foundation] " + path
			fields := []zap.Field{
				zap.String(conf.RequestIDField, reqID),
				zap.String("content_type", contentType),
				zap.Int("body_bytes", c.Writer.Size()),
				zap.Int("status_code", c.Writer.Status()),
				zap.String("latency", latency.String()),
				zap.String("client_ip", c.ClientIP()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
			}

			if len(c.Errors) > 0 {
				fields = append(fields,
					zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()))
				logger.Error(msg, fields...)
			} else {
				logger.Info(msg, fields...)
			}
		}
	}
}
