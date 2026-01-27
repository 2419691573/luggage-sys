package handlers

import (
	"net/http"
	"strconv"

	"luggage-sys2/internal/services"
	"luggage-sys2/internal/utils"

	"github.com/gin-gonic/gin"
)

type LuggageHandler struct {
	luggageService *services.LuggageService
}

func NewLuggageHandler() *LuggageHandler {
	return &LuggageHandler{
		luggageService: services.NewLuggageService(),
	}
}

func (h *LuggageHandler) CreateLuggage(c *gin.Context) {
	var req services.CreateLuggageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create luggage failed",
			"error":   "invalid request",
		})
		return
	}

	hotelID := utils.GetUintFromContext(c, "hotel_id")
	username := utils.GetStringFromContext(c, "username")
	if req.StaffName == "" {
		req.StaffName = username
	}

	luggage, code, err := h.luggageService.CreateLuggage(req, hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create luggage failed",
			"error":   err.Error(),
		})
		return
	}

	// 判断是多件模式还是单件模式
	if len(req.Items) > 0 {
		// 多件模式：返回所有创建的行李记录
		allLuggages, _ := h.luggageService.GetLuggageByCode(code, hotelID)
		items := make([]gin.H, 0, len(allLuggages))
		for _, l := range allLuggages {
			items = append(items, gin.H{
				"luggage_id":  l.ID,
				"storeroom_id": l.StoreroomID,
				"photo_url":   l.PhotoURL,
				"photo_urls":  l.PhotoURLs,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"message":       "create luggage success",
			"retrieval_code": code,
			"items":         items,
		})
	} else {
		// 单件模式：保持原有响应格式
		c.JSON(http.StatusOK, gin.H{
			"message":       "create luggage success",
			"luggage_id":    luggage.ID,
			"retrieval_code": code,
			"qrcode_url":    "/qr/" + code,
			"photo_url":     luggage.PhotoURL,
			"photo_urls":    luggage.PhotoURLs,
		})
	}
}

func (h *LuggageHandler) GetLuggageByCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "query luggage failed",
			"error":   "code is empty",
		})
		return
	}

	hotelID := utils.GetUintFromContext(c, "hotel_id")
	luggages, err := h.luggageService.GetLuggageByCode(code, hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "query luggage failed",
			"error":   err.Error(),
		})
		return
	}

	// 返回所有同取件码的行李记录
	items := make([]gin.H, 0, len(luggages))
	for _, luggage := range luggages {
		items = append(items, gin.H{
			"id":            luggage.ID,
			"guest_name":    luggage.GuestName,
			"contact_phone": luggage.ContactPhone,
			"storeroom_id":  luggage.StoreroomID,
			"retrieval_code": luggage.RetrievalCode,
			"status":        luggage.Status,
			"photo_url":     luggage.PhotoURL,
			"photo_urls":    luggage.PhotoURLs,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "query luggage success",
		"items":   items,
	})
}

func (h *LuggageHandler) CheckoutLuggage(c *gin.Context) {
	code := c.Param("id")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "checkout failed",
			"error":   "code is empty",
		})
		return
	}

	hotelID := utils.GetUintFromContext(c, "hotel_id")
	username := utils.GetStringFromContext(c, "username")

	luggageIDs, err := h.luggageService.CheckoutLuggage(code, username, hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "checkout failed",
			"error":   err.Error(),
		})
		return
	}

	// 根据文档，多件模式返回 retrieved_count 和 luggage_ids，单件模式返回 luggage_id
	if len(luggageIDs) > 1 {
		c.JSON(http.StatusOK, gin.H{
			"message":        "checkout success",
			"retrieval_code": code,
			"retrieved_count": len(luggageIDs),
			"luggage_ids":    luggageIDs,
			"luggage_id":     nil,
		})
	} else if len(luggageIDs) == 1 {
		c.JSON(http.StatusOK, gin.H{
			"message":        "checkout success",
			"retrieval_code": code,
			"retrieved_count": 1,
			"luggage_ids":    luggageIDs,
			"luggage_id":     luggageIDs[0],
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":        "checkout success",
			"retrieval_code": code,
			"retrieved_count": 0,
			"luggage_ids":    []uint{},
			"luggage_id":     nil,
		})
	}
}

func (h *LuggageHandler) GetGuestList(c *gin.Context) {
	hotelID := utils.GetUintFromContext(c, "hotel_id")
	if hotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing user info",
		})
		return
	}

	guestNames, err := h.luggageService.GetGuestList(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get checkout info failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get checkout info success",
		"items":   guestNames,
	})
}

func (h *LuggageHandler) GetLuggageByGuestName(c *gin.Context) {
	guestName := c.Query("guest_name")
	if guestName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list luggage failed",
			"error":   "guest_name is empty",
		})
		return
	}

	hotelID := utils.GetUintFromContext(c, "hotel_id")
	luggages, err := h.luggageService.GetLuggageByGuestName(guestName, hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list luggage failed",
			"error":   err.Error(),
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
			"photo_url":     luggage.PhotoURL,
			"photo_urls":    luggage.PhotoURLs,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list luggage success",
		"items":   items,
	})
}

func (h *LuggageHandler) UpdateLuggage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update luggage failed",
			"error":   "invalid luggage id",
		})
		return
	}

	var req services.UpdateLuggageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update luggage failed",
			"error":   "invalid request",
		})
		return
	}

	hotelID := utils.GetUintFromContext(c, "hotel_id")
	username := utils.GetStringFromContext(c, "username")

	if err := h.luggageService.UpdateLuggage(uint(id), req, hotelID, username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update luggage success",
	})
}
