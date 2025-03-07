package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stipochka/web_service/internal/models"
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

// @Summary Records
// @Tags records
// @Description records
// @Produce json
// @Success 200 {object} models.Record "record"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/records/:id [get]
func (h *Handler) getRecordById(c *gin.Context) {
	const op = "handler.getAllRecords"

	log := h.log.With(slog.String("op", op))
	idFromUrl, ok := c.Params.Get("id")
	if !ok {
		log.Error("failed to get id from url")

		newErrorResponse(c, http.StatusBadRequest, "not given id")

		return
	}
	mcuId, err := strconv.Atoi(idFromUrl)
	if err != nil {
		log.Error("failed to get mcuId", slog.Any("error", err))

		newErrorResponse(c, http.StatusBadRequest, "invalid id")

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

// @Summary Records
// @Tags records
// @Description records
// @Produce json
// @Success 200 {object} []models.Record "record"
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/records [get]
func (h *Handler) getAllRecords(c *gin.Context) {
	const op = "handler.getAllRecords"

	log := h.log.With(slog.String("op", op))

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
