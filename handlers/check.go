package handlers

import (
	"belajar-go/database"
	"belajar-go/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// CheckWebsites godoc
// @Summary      Cek status semua website (Manual Trigger)
// @Description  Memicu pengecekan status UP/DOWN secara paralel menggunakan Goroutines.
// @Tags         Monitor
// @Produce      json
// @Success      200  {array}   models.WebsiteStatus
// @Router       /check [post]
func CheckWebsites(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, url FROM websites")
	if err != nil {
		http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
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
		http.Error(w, "Belum ada website yang didaftarkan.", http.StatusBadRequest)
		return
	}

	resultsChan := make(chan models.WebsiteStatus, len(websites))
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

			_, dbErr := database.DB.Exec(
				"INSERT INTO checks (website_id, status) VALUES ($1, $2)",
				targetWeb.ID, status,
			)
			if dbErr != nil {
				fmt.Println("[Worker Error] Gagal mencatat log untuk", targetWeb.URL)
			}

			resultsChan <- models.WebsiteStatus{URL: targetWeb.URL, Status: status}
		}(web)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var checkResults []models.WebsiteStatus
	for res := range resultsChan {
		checkResults = append(checkResults, res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkResults)
}
