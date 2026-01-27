package handlers

import (
	"net/http"
	"strconv"

	"luggage-sys2/internal/services"
	"luggage-sys2/internal/utils"

	"github.com/gin-gonic/gin"
)

type StoreroomHandler struct {
	storeroomService *services.StoreroomService
}

func NewStoreroomHandler() *StoreroomHandler {
	return &StoreroomHandler{
		storeroomService: services.NewStoreroomService(),
	}
}

func (h *StoreroomHandler) ListStorerooms(c *gin.Context) {
	hotelID := utils.GetUintFromContext(c, "hotel_id")
	if hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list storerooms failed",
			"error":   "hotel_id is missing",
		})
		return
	}

	storerooms, err := h.storeroomService.ListStorerooms(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list storerooms failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list storerooms success",
		"items":   storerooms,
	})
}

func (h *StoreroomHandler) CreateStoreroom(c *gin.Context) {
	var req services.CreateStoreroomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create storeroom failed",
			"error":   "invalid request",
		})
		return
	}

	hotelID := utils.GetUintFromContext(c, "hotel_id")
	storeroom, err := h.storeroomService.CreateStoreroom(req, hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create storeroom failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create storeroom success",
		"item":     storeroom,
	})
}

func (h *StoreroomHandler) UpdateStoreroom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update storeroom status failed",
			"error":   "invalid storeroom id",
		})
		return
	}

	var req services.UpdateStoreroomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update storeroom status failed",
			"error":   "invalid request",
		})
		return
	}

	hotelID := utils.GetUintFromContext(c, "hotel_id")
	if err := h.storeroomService.UpdateStoreroom(uint(id), req, hotelID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update storeroom status failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update storeroom status success",
	})
}

func (h *StoreroomHandler) GetStoreroomOrders(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid storeroom id",
		})
		return
	}

	hotelID := utils.GetUintFromContext(c, "hotel_id")
	status := c.Query("status")

	luggages, err := h.storeroomService.GetStoreroomOrders(uint(id), hotelID, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	items := make([]gin.H, 0, len(luggages))
	for _, luggage := range luggages {
		items = append(items, gin.H{
			"id":            luggage.ID,
			"guest_name":    luggage.GuestName,
			"retrieval_code": luggage.RetrievalCode,
			"status":        luggage.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list luggage success",
		"items":   items,
	})
}
