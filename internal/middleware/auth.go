package middleware

import (
	"net/http"
	"strings"

	"luggage-sys2/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "missing authorization header",
			})
			c.Abort()
			return
		}

		// 检查 Bearer 前缀（兼容多空格）
		fields := strings.Fields(authHeader)
		if len(fields) != 2 || fields[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := fields[1]
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", int(claims.UserID))
		c.Set("username", claims.Username)
		c.Set("hotel_id", int(claims.HotelID))

		c.Next()
	}
}
