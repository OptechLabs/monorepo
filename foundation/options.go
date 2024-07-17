package foundation

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Options struct {
	Environment                 string                      `long:"environment" description:"environment to run in" default:"development"`
	HTTPPort                    string                      `long:"proxy-port" description:"port which http server listens on" default:"8080"`
	StartHTTPServer             bool                        `long:"start-http-server" description:"run the http server" default:"false"`
	GRPCPort                    string                      `long:"proxy-port" description:"port which grpc server listens on" default:"8081"`
	GRPCUnaryInterceptor        grpc.UnaryServerInterceptor `long:"-" description:"grpc unary interceptors"`
	StartGRPCServer             bool                        `long:"start-grpc-server" description:"run the grpc server" default:"false"`
	Logger                      Logger                      `long:"-" description:"logger"`
	WriteTimeout                time.Duration               `long:"write-timeout" description:"http server write timeout" default:"15s"`
	ReadTimeout                 time.Duration               `long:"read-timeout" description:"http server read timeout" default:"15s"`
	IdleTimeout                 time.Duration               `long:"idle-timeout" description:"http server idle timeout" default:"60s"`
	ShutdownWait                time.Duration               `long:"shutdown-wait" description:"time to wait for server to shutdown" default:"30s"`
	StopOnProcessorStartFailure bool                        `long:"stop-on-processor-start-failure" description:"stop the server if a processor fails to start"`
}

func (o Options) ValuesOrDefaults() Options {
	if o.Environment == "" {
		o.Environment = "development"
	}
	if o.HTTPPort == "" {
		o.HTTPPort = "8080"
	}
	if o.GRPCPort == "" {
		o.GRPCPort = "8081"
	}
	if o.Logger == nil {
		o.Logger, _ = NewDefaultLogger(o.Environment)
	}
	if o.WriteTimeout == 0 {
		o.WriteTimeout = 15 * time.Second
	}
	if o.ReadTimeout == 0 {
		o.ReadTimeout = 15 * time.Second
	}
	if o.IdleTimeout == 0 {
		o.IdleTimeout = 60 * time.Second
	}
	if o.ShutdownWait == 0 {
		o.ShutdownWait = 30 * time.Second
	}
	return o
}

const (
	Development = "development"
	Test        = "test"
	Staging     = "staging"
	Sandbox     = "sandbox"
	Integration = "integration"
	Production  = "production"
)

func (o Options) Mode() string {
	switch strings.ToLower(o.Environment) {
	case Development:
		return gin.DebugMode
	case Test:
		return gin.TestMode
	case Staging, Integration, Sandbox, Production:
		return gin.ReleaseMode
	}
	return gin.ReleaseMode
}
