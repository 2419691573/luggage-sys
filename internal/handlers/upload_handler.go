package handlers

import (
	"log"
	"net/http"
	"strings"

	"luggage-sys2/internal/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService *services.UploadService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		uploadService: services.NewUploadService(),
	}
}

// UploadImage handles:
//   POST /api/upload
// Content-Type: multipart/form-data
// Form field:
//   file: image file (jpg/png/webp)
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 记录请求信息用于调试
	contentType := c.GetHeader("Content-Type")
	log.Printf("[Upload] Request Content-Type: %s", contentType)
	log.Printf("[Upload] Request Method: %s", c.Request.Method)
	
	// 检查 Content-Type 是否为 multipart/form-data
	if !strings.Contains(strings.ToLower(contentType), "multipart/form-data") {
		log.Printf("[Upload] Invalid Content-Type: %s, expected multipart/form-data", contentType)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "upload failed",
			"error":   "invalid content type",
			"detail":  "Content-Type must be 'multipart/form-data', but got: " + contentType,
			"hint":    "Please use FormData to send the file. Do not set Content-Type header manually, let the browser set it automatically.",
		})
		return
	}
	
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("[Upload] FormFile error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "upload failed",
			"error":   "missing or invalid file field",
			"detail":  err.Error(),
			"hint":    "Make sure the form field name is 'file' and the file is sent as multipart/form-data",
		})
		return
	}

	if fileHeader == nil {
		log.Printf("[Upload] fileHeader is nil")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "upload failed",
			"error":   "file header is nil",
		})
		return
	}

	log.Printf("[Upload] Received file: %s, size: %d bytes", fileHeader.Filename, fileHeader.Size)

	// 5MB
	const maxBytes = 5 * 1024 * 1024
	res, err := h.uploadService.SaveImageFile(fileHeader, maxBytes)
	if err != nil {
		log.Printf("[Upload] SaveImageFile error: %v", err)
		msg := "upload failed"
		code := http.StatusBadRequest
		if err == services.ErrFileTooLarge {
			c.JSON(code, gin.H{"message": msg, "error": "file too large"})
			return
		}
		if err == services.ErrInvalidFileType {
			c.JSON(code, gin.H{"message": msg, "error": "invalid file type"})
			return
		}
		if err == services.ErrMissingFileField {
			c.JSON(code, gin.H{"message": msg, "error": "missing file field"})
			return
		}
		// 其他错误可能是内部错误，返回 500
		log.Printf("[Upload] Internal error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": msg,
			"error":   "internal server error",
			"detail":  err.Error(),
		})
		return
	}

	// Build absolute URL for convenience
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := strings.TrimSpace(c.Request.Host)
	absoluteURL := res.RelativeURL
	if host != "" {
		absoluteURL = scheme + "://" + host + res.RelativeURL
	}

	log.Printf("[Upload] Upload success: %s", res.RelativeURL)
	c.JSON(http.StatusOK, gin.H{
		"message":       "upload success",
		"url":           absoluteURL,
		"relative_url":  res.RelativeURL,
		"content_type":  res.ContentType,
		"size":          res.Size,
		"file_name":     res.FileName,
		"max_size_byte": maxBytes,
	})
}

