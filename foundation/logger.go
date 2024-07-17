package foundation

import (
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: REplace with slog golang standard library

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
}

func NewDefaultLogger(environment string) (Logger, error) {
	config := zap.Config{}

	if environment == "" {
		environment = "development"
	}

	switch environment {
	case Development, Test:
		config = zap.NewDevelopmentConfig()
	case Staging, Sandbox, Integration, Production:
		config = zap.NewProductionConfig()
	}

	config.Encoding = "json"
	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		lvl := l.CapitalString()
		if lvl == "WARN" {
			lvl = "WARNING"
		}
		enc.AppendString(lvl + ":")
	}
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return config.Build()
}

func LogExecutionTime(l Logger, methodOverride string, fn func()) {
	getCallerName := func() (method string) {
		pc, _, _, _ := runtime.Caller(3)
		fn := runtime.FuncForPC(pc)
		arr := strings.Split(fn.Name(), "/")
		method = arr[len(arr)-1]
		return
	}
	startTime := time.Now()
	fn()
	method := methodOverride
	if method == "" {
		method = getCallerName()
	}
	l.Info("otos."+method, zap.Duration("executionTime", time.Duration(time.Since(startTime).Milliseconds())))
}
