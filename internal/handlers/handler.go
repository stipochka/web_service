package handler

import (
	"log/slog"

	_ "github.com/stipochka/web_service/docs"
	"github.com/stipochka/web_service/internal/service"
	swaggerfiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"

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

		router.GET("/ws/records", h.getAllRecordsWebSocket)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
