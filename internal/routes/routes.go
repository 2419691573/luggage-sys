package routes

import (
	"luggage-sys2/internal/handlers"
	"luggage-sys2/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// 设置 multipart 表单大小限制为 10MB（大于文件大小限制 5MB）
	r.MaxMultipartMemory = 10 << 20 // 10MB

	// 添加 CORS 中间件（必须在所有路由之前）
	r.Use(middleware.CORSMiddleware())

	// 静态资源访问：本地存储
	// 访问示例：GET /uploads/2026/01/xxx.jpg
	r.Static("/uploads", "./uploads")

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 认证处理器
	authHandler := handlers.NewAuthHandler()

	// API路由组
	api := r.Group("/api")
	{
		// 登录（不需要认证）
		api.POST("/login", authHandler.Login)

		// 需要认证的路由
		api.Use(middleware.AuthMiddleware())
		{
			// 图片上传（不影响行李接口结构：上传后把返回的 url 写入 photo_url）
			uploadHandler := handlers.NewUploadHandler()
			api.POST("/upload", uploadHandler.UploadImage)

			// 行李相关路由
			luggageHandler := handlers.NewLuggageHandler()
			api.POST("/luggage", luggageHandler.CreateLuggage)
			api.GET("/luggage/by_code", luggageHandler.GetLuggageByCode)
			api.POST("/luggage/:id/checkout", luggageHandler.CheckoutLuggage)
			api.GET("/luggage/:id/checkout", luggageHandler.GetGuestList)
			api.GET("/luggage/list/by_guest_name", luggageHandler.GetLuggageByGuestName)
			api.PUT("/luggage/:id", luggageHandler.UpdateLuggage)

			// 寄存室相关路由
			storeroomHandler := handlers.NewStoreroomHandler()
			api.GET("/luggage/storerooms", storeroomHandler.ListStorerooms)
			api.POST("/luggage/storerooms", storeroomHandler.CreateStoreroom)
			api.PUT("/luggage/storerooms/:id", storeroomHandler.UpdateStoreroom)
			api.GET("/luggage/storerooms/:id/orders", storeroomHandler.GetStoreroomOrders)

			// 日志相关路由
			logHandler := handlers.NewLogHandler()
			api.GET("/luggage/logs/stored", logHandler.GetStoredLogs)
			api.GET("/luggage/logs/updated", logHandler.GetUpdatedLogs)
			api.GET("/luggage/logs/retrieved", logHandler.GetRetrievedLogs)
		}
	}

	return r
}
