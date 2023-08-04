package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/antonmisa/1cctl/config"
	v1 "github.com/antonmisa/1cctl/internal/controller/v1"
	"github.com/antonmisa/1cctl/internal/usecase"
	uccache "github.com/antonmisa/1cctl/internal/usecase/cache"
	ucpipe "github.com/antonmisa/1cctl/internal/usecase/pipe"
	"github.com/antonmisa/1cctl/pkg/cache"
	"github.com/antonmisa/1cctl/pkg/httpserver"
	"github.com/antonmisa/1cctl/pkg/logger"
	"github.com/antonmisa/1cctl/pkg/pipe"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	c, err := cache.New(cfg.Redis.Addr, cfg.Redis.TTL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - RunHTTP - cache.New: %w", err))
	}

	p, err := pipe.New(cfg.App.PathToRAC)
	if err != nil {
		l.Fatal(fmt.Errorf("app - RunHTTP - pipe.New: %w", err))
	}

	// Use case
	ctrlUseCase := usecase.New(
		uccache.New(c),
		ucpipe.New(p),
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, ctrlUseCase)
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
