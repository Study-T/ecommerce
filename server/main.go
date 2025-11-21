package main

import (
	"ecommerce/router"
	"ecommerce/service"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Gin 引擎
	r := gin.Default()

	// 设置上传文件大小限制 (32MB)
	r.MaxMultipartMemory = 32 << 20

	// 获取上传目录路径
	uploadPath := filepath.Join(".", "uploads")

	// 初始化上传服务（支持空文件夹上传）
	uploadService := service.NewUploadService(uploadPath)

	// 注册路由
	uploadRouter := router.NewUploadRouter(uploadService)
	uploadRouter.RegisterRoutes(r)

	// 添加健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "服务器运行正常，支持空文件夹上传",
		})
	})

	// 启动服务器
	log.Println("服务器启动在 :8080")
	log.Println("支持功能: 文件上传、文件夹上传（包括空文件夹）")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
