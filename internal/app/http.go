package app

import (
	"fmt"
	"github.com/antonmisa/1cctl/internal/app/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"os"
	"os/signal"
	"syscall"

	"github.com/antonmisa/1cctl/config"
	v1 "github.com/antonmisa/1cctl/internal/controller/http/v1"
	"github.com/antonmisa/1cctl/internal/usecase"
	ucbackup "github.com/antonmisa/1cctl/internal/usecase/backup"
	uccache "github.com/antonmisa/1cctl/internal/usecase/cache"
	ucpipe "github.com/antonmisa/1cctl/internal/usecase/pipe"
	"github.com/antonmisa/1cctl/pkg/cache"
	"github.com/antonmisa/1cctl/pkg/httpserver"
	"github.com/antonmisa/1cctl/pkg/logger"
	"github.com/antonmisa/1cctl/pkg/pipe"
)

func Run(cfg *config.Config) {
	l, err := logger.New(cfg.Log.Path, cfg.Log.Level)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - logger.New: %w", err))
	}

	c, err := cache.New(cfg.Cache.TTL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - cache.New: %w", err))
	}

	p, err := pipe.New(cfg.App.PathToRAC)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - pipe.New: %w", err))
	}

	cb, err := ucbackup.New(cfg.App.PathTo1C)

	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - ucbackup.New: %w", err))
	}

	// Use case
	ctrlUseCase := usecase.New(
		uccache.New(c),
		ucpipe.New(p),
		cb,
	)

	// Trace start
	tp, err := tracing.JaegerTraceProvider(cfg.Trace.Endpoint)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - tracing.JaegerTraceProvider: %w", err))
	}

	if cfg.Trace.Enable {
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, ctrlUseCase, tp.Tracer("1ctrl_main_trace"))
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - RunHTTP - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - RunHTTP - httpServer.Shutdown: %w", err))
	}
}
