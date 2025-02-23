package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stipochka/web_service/internal/models"
)

const (
	mcuCtx = "mcuId"
)

type errorResponse struct {
	Message string `json:"message"`
}

type getAllRecordsResponse struct {
	Data []models.Record `json:"data"`
}

func newErrorResponse(c *gin.Context, statusCode int, errMessage string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: errMessage})
}

func getMcuId(c *gin.Context) (int, error) {
	id, ok := c.Get(mcuCtx)
	if !ok {
		return 0, errors.New("failed to get mcu Id")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("invalid id")
	}

	return idInt, nil

}

func (h *Handler) getRecordById(c *gin.Context) {
	const op = "handler.getAllRecords"

	log := h.log.With(op)

	mcuId, err := getMcuId(c)
	if err != nil {
		log.Error("failed to get mcuId", slog.Any("error", err))

		newErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	record, err := h.services.GetRecordById(context.Background(), mcuId)
	if err != nil {
		log.Error("failed to get record", slog.Any("error", err.Error()))

		newErrorResponse(c, http.StatusInternalServerError, "internal error")

		return
	}

	log.Info("received record with mcuId", slog.Int("mcuID", mcuId))
	c.JSON(http.StatusOK, record)
}

func (h *Handler) getAllRecords(c *gin.Context) {
	const op = "handler.getAllRecords"

	log := h.log.With(op)

	records, err := h.services.GetAllRecords(context.Background())
	if err != nil {
		log.Error("failed to get records", slog.Any("error", err))

		newErrorResponse(c, http.StatusInternalServerError, "Internal error")

		return
	}

	log.Info("received records")
	c.JSON(http.StatusOK, getAllRecordsResponse{
		Data: records,
	})
}
