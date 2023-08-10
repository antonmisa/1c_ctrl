// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-contrib/pprof"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/antonmisa/1cctl/docs"
	"github.com/antonmisa/1cctl/internal/controller/http/v1/middleware/commonqueryparams"
	mwlogger "github.com/antonmisa/1cctl/internal/controller/http/v1/middleware/logger"
	"github.com/antonmisa/1cctl/internal/usecase"
	"github.com/antonmisa/1cctl/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       1C cluster control service
// @description Using a 1C cluster control service over http
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.Ctrl, tr trace.Tracer) {
	// Options
	handler.Use(mwlogger.Logger(l))
	handler.Use(gin.Recovery())
	handler.Use(otelgin.Middleware("1ctrl-service"))

	// Pprof
	pprof.Register(handler)

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		h.Use(commonqueryparams.UseCommonQueryParams(l))

		newCtrlRoutes(h, t, l, tr)
	}
}
