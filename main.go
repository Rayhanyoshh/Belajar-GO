package main

import (
	"belajar-go/database"
	_ "belajar-go/docs"
	"belajar-go/handlers"
	"belajar-go/middlewares"
	"belajar-go/workers"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		err := database.DB.Ping()
		dbStatus := "connected"
		if err != nil {
			dbStatus = "disconnected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "gotracker-main",
			"db":        dbStatus,
			"timestamp": time.Now().UTC(),
		})
	})

	// Rute publik
	r.GET("/websites", handlers.GetWebsites)
	r.GET("/stats", handlers.GetStats)
	
	// Route khusus untuk membuka halaman dokumentasi visual Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// RUTE RAHASIA! (Dibungkus dengan Middleware Satpam kita)
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/websites", handlers.AddWebsite)
		protected.GET("/websites/:id/history", handlers.GetWebsiteHistory)
		protected.DELETE("/websites/:id", handlers.DeleteWebsite)
		protected.PUT("/websites/:id", handlers.UpdateWebsite)
		protected.POST("/check", handlers.CheckWebsites)
	}

	slog.Info("Server Tracker mulai berjalan", "port", "8080")
	slog.Info("Dokumentasi API tersedia di /swagger/index.html")
	
	// Mendukung port dinamis untuk deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error saat menjalankan server", "error", err)
		}
	}()

	// Menunggu sinyal interrupt untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Sinyal interrupt diterima, mematikan server Tracker...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server terpaksa mati", "error", err)
	}
	
	slog.Info("Server Tracker telah berhenti dengan aman.")
}
