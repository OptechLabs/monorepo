package foundation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Foundation struct {
	Environment                 string
	ShutdownWaitTime            time.Duration
	HTTPRouter                  *gin.Engine
	HTTPServer                  *http.Server
	GRPCServer                  *grpc.Server
	StopOnProcessorStartFailure bool
	ShutdownWait                time.Duration
	Logger                      Logger
	ctx                         context.Context
	processors                  []Processor
	startopOpts                 Options
}

// New creates a new Foundation instance. It does not start it up.
//
// Default Settings:
// - Environment: development
// - Port: 3000
// - ShutdownWait: 30 seconds
// - WriteTimeout: 15 seconds
// - ReadTimeout: 15 seconds
// - IdleTimeout: 60 seconds
func New(opts Options) *Foundation {
	opts = opts.ValuesOrDefaults()
	gin.SetMode(opts.Mode())

	router := gin.New()
	return &Foundation{
		Environment: opts.Environment,
		HTTPRouter:  router,
		HTTPServer: func(runHTTP bool, router *gin.Engine, opts Options) *http.Server {
			if !runHTTP {
				return nil
			}
			return &http.Server{
				Addr:           fmt.Sprintf("0.0.0.0:%s", opts.HTTPPort),
				WriteTimeout:   opts.WriteTimeout,
				ReadTimeout:    opts.ReadTimeout,
				IdleTimeout:    opts.IdleTimeout,
				Handler:        router,
				MaxHeaderBytes: 1 << 20,
			}
		}(opts.StartHTTPServer, router, opts),
		GRPCServer: func(runGRPC bool) *grpc.Server {
			if !runGRPC {
				return nil
			}
			if opts.GRPCUnaryInterceptor != nil {
				return grpc.NewServer(withServerUnaryInterceptors(opts.GRPCUnaryInterceptor))
			}
			return grpc.NewServer()
		}(opts.StartGRPCServer),
		Logger: func(optLog Logger) Logger {
			if optLog == nil {
				l, err := NewDefaultLogger(opts.Environment)
				if err != nil {
					log.Fatal(err)
				}
				return l
			}
			return optLog
		}(opts.Logger),
		startopOpts: opts,
	}
}

func withServerUnaryInterceptors(interceptor grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.UnaryInterceptor(interceptor)
}

// Serve starts the foundation server and your app.
// func (f *Foundation) Serve(quit <-chan os.Signal) error {
func (f *Foundation) RunWithContext(ctx context.Context, stop context.CancelFunc) error {
	if errs := f.StartProcessors(); len(errs) > 0 {
		if f.StopOnProcessorStartFailure {
			return errors.New("[foundation] ERROR: foundation failed to start one or more attached processors. StopOnProcessorStartFailure setting is true")
		}
	}

	if f.GRPCServer != nil {
		go func(port string) {
			f.Logger.Info("grpc server starting", zap.String("port", port))
			listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
			if err != nil {
				f.Logger.Fatal("grpc listener failed to start", zap.Error(err), zap.String("port", port))
			}

			if err := f.GRPCServer.Serve(listener); err != nil {
				f.Logger.Fatal("grpc server failed to start", zap.Error(err), zap.String("port", port))
			}
			f.Logger.Info("grpc server started", zap.String("port", port))
		}(f.startopOpts.GRPCPort)
	}

	if f.HTTPServer != nil {
		go func() {
			f.Logger.Info("http server starting", zap.String("tcpAddress", f.HTTPServer.Addr))
			if err := f.HTTPServer.ListenAndServe(); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					stop()
				}
				f.Logger.Fatal("http server failed to start", zap.Error(err), zap.String("tcpAddress", f.HTTPServer.Addr))
			}
			f.Logger.Info("http server started", zap.String("tcpAddress", f.HTTPServer.Addr))
		}()
	}

	// Block until we receive our signal.
	<-ctx.Done()
	f.Logger.Info("shutting down server")

	var wg sync.WaitGroup
	if errs := f.StopProcessors(&wg); len(errs) > 0 {
		f.Logger.Info("foundation failed to gracefully shutdown one or more attached processors", zap.Errors("errors", errs))
	}
	wg.Wait()

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), f.ShutdownWait)
	defer cancel()

	if f.GRPCServer != nil {
		f.Logger.Info("shutting down grpc server")
		f.GRPCServer.GracefulStop()
	}
	if f.HTTPServer != nil {
		f.Logger.Info("shutting down http server")
		if err := f.HTTPServer.Shutdown(ctx); err != nil {
			return err
		}
	}

	f.Logger.Info("foundation stopped")
	return nil
}

func (f *Foundation) Run() error {
	return f.RunWithContext(ContextWithCancel())
}

func ContextWithCancel() (context.Context, context.CancelFunc) {
	ctx, stop := context.WithCancel(context.Background())
	CancelOnSignal(stop)
	return ctx, stop
}

func CancelOnSignal(stop context.CancelFunc) {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-stopSignal
		log.Printf("[foundation] stop signal received: %+v", s)
		stop()
	}()
}
