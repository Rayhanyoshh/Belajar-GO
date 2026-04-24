package main

import (
	"belajar-go/database"
	_ "belajar-go/docs"
	"belajar-go/handlers"
	"belajar-go/middlewares"
	"belajar-go/workers"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           GoTracker API
// @version         1.0
// @description     Sistem Pemantau Status Website Berbasis Goroutines
// @host            localhost:8080
// @BasePath        /
func main() {
	database.ConnectDB()

	workers.StartBackgroundChecker()

	r := gin.Default()

	// 3. Konfigurasi CORS (Gin)
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rute publik
	r.GET("/websites", handlers.GetWebsites)
	
	// Route khusus untuk membuka halaman dokumentasi visual Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// RUTE RAHASIA! (Dibungkus dengan Middleware Satpam kita)
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/websites", handlers.AddWebsite)
		protected.POST("/check", handlers.CheckWebsites)
	}

	fmt.Println("=======================================")
	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	fmt.Println("📖 Dokumentasi API: http://localhost:8080/swagger/index.html")
	fmt.Println("=======================================")
	
	r.Run(":8080")
}
