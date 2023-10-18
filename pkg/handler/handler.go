package handler

import (
	"github.com/MarselBissengaliyev/advertising/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	_ "github.com/MarselBissengaliyev/advertising/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	router.GET("/swagger/*any", echoSwagger.WrapHandler)

	api := router.Group("/api")
	adverts := api.Group("/adverts")
	adverts.GET("/", h.getAllAdverts)
	adverts.GET("/:id", h.getAdvertById)
	adverts.POST("/", h.createAdvert)

	return router
}
