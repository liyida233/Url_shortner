package main

import (
	"URL_shortener/handler"
	"URL_shortener/store"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	mysqlDB := store.InitMySQL()
	redisClient := store.InitRedis()

	storeService := store.NewHybridStore(
		store.NewRedisStore(redisClient),
		store.NewMySQLStore(mysqlDB),
	)

	// 初始化 handler
	h := handler.NewHandler(storeService)

	// 设置路由
	r := gin.Default()

	// 添加这行用于服务静态前端资源
	r.Static("/static", "./static")

	// 示例：访问 http://localhost:9808/static/index.html
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.POST("/create-short-url", h.CreateShortUrl)
	r.GET("/:shortUrl", h.HandleShortUrlRedirect)

	// 启动服务
	if err := r.Run("0.0.0.0:9808"); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
