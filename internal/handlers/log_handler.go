package handlers

import (
	"net/http"

	"luggage-sys2/internal/services"
	"luggage-sys2/internal/utils"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logService *services.LogService
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		logService: services.NewLogService(),
	}
}

func (h *LogHandler) GetStoredLogs(c *gin.Context) {
	hotelID := utils.GetUintFromContext(c, "hotel_id")
	if hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list logs failed",
			"error":   "hotel_id is missing",
		})
		return
	}

	logs, err := h.logService.GetStoredLogs(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list logs failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list logs success",
		"items":   logs,
	})
}

func (h *LogHandler) GetUpdatedLogs(c *gin.Context) {
	hotelID := utils.GetUintFromContext(c, "hotel_id")
	if hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list logs failed",
			"error":   "hotel_id is missing",
		})
		return
	}

	logs, err := h.logService.GetUpdatedLogs(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list logs failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list logs success",
		"items":   logs,
	})
}

func (h *LogHandler) GetRetrievedLogs(c *gin.Context) {
	hotelID := utils.GetUintFromContext(c, "hotel_id")
	if hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list logs failed",
			"error":   "hotel_id is missing",
		})
		return
	}

	logs, err := h.logService.GetRetrievedLogs(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list logs failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list logs success",
		"items":   logs,
	})
}
