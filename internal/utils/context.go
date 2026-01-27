package utils

import (
	"github.com/gin-gonic/gin"
)

// GetUintFromContext 从上下文中获取uint值
func GetUintFromContext(c *gin.Context, key string) uint {
	value, exists := c.Get(key)
	if !exists {
		return 0
	}
	if intValue, ok := value.(int); ok {
		return uint(intValue)
	}
	return 0
}

// GetStringFromContext 从上下文中获取string值
func GetStringFromContext(c *gin.Context, key string) string {
	value, exists := c.Get(key)
	if !exists {
		return ""
	}
	if strValue, ok := value.(string); ok {
		return strValue
	}
	return ""
}
