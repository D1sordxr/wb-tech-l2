package middleware

import (
	"time"
	appPorts "wb-tech-l2/18/calendar/internal/domain/app/ports"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(log appPorts.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start)
		log.Info("request finished",
			"time_elapsed", elapsed,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
	}
}
