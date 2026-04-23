package main

import (
	"belajar-go/database"
	_ "belajar-go/docs" // Wajib: Import hasil generate dokumen swagger nanti
	"belajar-go/handlers"
	"belajar-go/middlewares"
	"belajar-go/workers"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           GoTracker API
// @version         1.0
// @description     Sistem Pemantau Status Website Berbasis Goroutines
// @host            localhost:8080
// @BasePath        /
func main() {
	database.ConnectDB()

	workers.StartBackgroundChecker()

	mux := http.NewServeMux()
	
	// Rute publik (Siapa saja boleh melihat daftar)
	mux.HandleFunc("GET /websites", handlers.GetWebsites)
	
	// RUTE RAHASIA! (Dibungkus dengan Middleware Satpam kita)
	mux.HandleFunc("POST /websites", middlewares.AuthMiddleware(handlers.AddWebsite))
	mux.HandleFunc("POST /check", middlewares.AuthMiddleware(handlers.CheckWebsites))

	// Route khusus untuk membuka halaman dokumentasi visual Swagger UI
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Arahkan ke file spesifikasi swagger
	))

	// Bungkus seluruh router kita (Mux) dengan jubah pelindung CORS!
	handler := middlewares.CORSMiddleware(mux)

	fmt.Println("=======================================")
	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	fmt.Println("📖 Dokumentasi API: http://localhost:8080/swagger/index.html")
	fmt.Println("=======================================")
	
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println("Server gagal berjalan:", err)
	}
}
