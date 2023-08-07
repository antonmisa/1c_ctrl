package logger

import (
	"time"

	"github.com/antonmisa/1cctl/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		// Stop timer
		latency := time.Since(start)

		if raw != "" {
			path = path + "?" + raw
		}

		l.Info("%s   %s   %d   %s   %dmsec", c.Request.Method, path, c.Writer.Status(), c.ClientIP(), latency.Milliseconds())
	}
}
