package workers

import (
	"belajar-go/database"
	"belajar-go/models"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// StartBackgroundChecker akan dipanggil SATU KALI di main.go saat server menyala
func StartBackgroundChecker() {
	// time.Ticker adalah "jam weker" yang akan berdetak setiap interval tertentu
	// Untuk kebutuhan belajar, kita atur ke 15 detik (Di industri biasanya 1-5 menit)
	ticker := time.NewTicker(15 * time.Second)
	
	// Kita bungkus ke dalam Goroutine agar fungsi ini berjalan terus di background
	// tanpa memblokir server utama
	go func() {
		// Looping tak terbatas yang akan ter-trigger setiap kali ticker berdetak
		for range ticker.C {
			slog.Info("Memulai pengecekan otomatis", "source", "background_worker")
			runCheck()
		}
	}()
}

// Logika inti pekerja. Mirip sekali dengan yang ada di handlers/check.go
func runCheck() {
	rows, err := database.DB.Query("SELECT id, url FROM websites")
	if err != nil {
		slog.Error("Gagal mengambil data dari database", "error", err, "source", "background_worker")
		return
	}
	defer rows.Close()

	var websites []models.Website
	for rows.Next() {
		var web models.Website
		if err := rows.Scan(&web.ID, &web.URL); err == nil {
			websites = append(websites, web)
		}
	}

	if len(websites) == 0 {
		return // Tidak ada yang perlu dicek
	}

	var wg sync.WaitGroup

	for _, web := range websites {
		wg.Add(1)
		
		go func(targetWeb models.Website) {
			defer wg.Done()
			
			start := time.Now()
			resp, httpErr := http.Get(targetWeb.URL)
			latency := time.Since(start).Milliseconds()

			status := "UP"
			if httpErr != nil || resp.StatusCode >= 400 {
				status = "DOWN"
			}

			// Simpan hasil ke database
			_, dbErr := database.DB.Exec(
				"INSERT INTO checks (website_id, status, response_time_ms) VALUES ($1, $2, $3)",
				targetWeb.ID, status, latency,
			)
			if dbErr != nil {
				slog.Error("Gagal mencatat log", "url", targetWeb.URL, "error", dbErr, "source", "background_worker")
			} else {
				slog.Info("Check selesai", "url", targetWeb.URL, "status", status, "latency_ms", latency, "source", "background_worker")
			}
		}(web)
	}
	
	wg.Wait() // Tunggu sampai semua website selesai dicek
}
