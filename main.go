package main

import (
	"fmt"
	"log"

	"luggage-sys2/internal/config"
	"luggage-sys2/internal/database"
	"luggage-sys2/internal/routes"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	database.Init()

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	addr := config.Port

	fmt.Printf("Listening and serving HTTP on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

