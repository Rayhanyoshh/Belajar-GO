package workers

import (
	"belajar-go/database"
	"belajar-go/models"
	"fmt"
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
			fmt.Println("\n⏳ [Background Worker] Memulai pengecekan otomatis...")
			runCheck()
		}
	}()
}

// Logika inti pekerja. Mirip sekali dengan yang ada di handlers/check.go
func runCheck() {
	rows, err := database.DB.Query("SELECT id, url FROM websites")
	if err != nil {
		fmt.Println("[Worker Error] Gagal mengambil data dari database:", err)
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
			
			resp, httpErr := http.Get(targetWeb.URL)
			status := "UP"
			if httpErr != nil || resp.StatusCode >= 400 {
				status = "DOWN"
			}

			// Simpan hasil ke database
			_, dbErr := database.DB.Exec(
				"INSERT INTO checks (website_id, status) VALUES ($1, $2)",
				targetWeb.ID, status,
			)
			if dbErr != nil {
				fmt.Println("[Worker Error] Gagal mencatat log untuk", targetWeb.URL)
			} else {
				// Print ke terminal agar kita tahu robotnya sedang bekerja
				fmt.Printf("✅ [Worker] %s -> %s\n", targetWeb.URL, status)
			}
		}(web)
	}
	
	wg.Wait() // Tunggu sampai semua website selesai dicek
}
