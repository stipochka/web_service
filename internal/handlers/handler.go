package handler

import (
	"log/slog"

	"github.com/stipochka/web_service/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	log      *slog.Logger
	services *service.Service
}

func NewHandler(log *slog.Logger, service *service.Service) *Handler {
	return &Handler{
		log:      log,
		services: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		records := api.Group("/records")
		{
			records.GET("/", h.getAllRecords)
			records.GET("/:id", h.getRecordById)
		}
	}

	return router
}
