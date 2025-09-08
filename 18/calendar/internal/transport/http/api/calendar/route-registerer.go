package calendar

import (
	"wb-tech-l2/18/calendar/internal/transport/http/api/calendar/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RouteRegisterer struct {
	handlers    *handler.Handlers
	middlewares []gin.HandlerFunc
}

func NewRouteRegisterer(
	handlers *handler.Handlers,
	middlewares ...gin.HandlerFunc,
) *RouteRegisterer {
	return &RouteRegisterer{
		handlers:    handlers,
		middlewares: middlewares,
	}
}

func (r *RouteRegisterer) RegisterRoutes(router gin.IRouter) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})
	router.GET("/get-personal-id", func(c *gin.Context) {
		c.JSON(200, gin.H{"id": uuid.New()})
	})

	handler.RegisterHandlers(
		router.Group("/calendar"),
		handler.NewStrictHandler(
			r.handlers,
			[]handler.StrictMiddlewareFunc{},
		),
	)
}
