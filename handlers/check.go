package handlers

import (
	"belajar-go/database"
	"belajar-go/models"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CheckWebsites godoc
// @Summary      Cek status semua website (Manual Trigger)
// @Description  Memicu pengecekan status UP/DOWN secara paralel menggunakan Goroutines.
// @Tags         Monitor
// @Produce      json
// @Success      200  {array}   models.WebsiteStatus
// @Router       /check [post]
func CheckWebsites(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, url FROM websites")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Belum ada website yang didaftarkan."})
		return
	}

	resultsChan := make(chan models.WebsiteStatus, len(websites))
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

			_, dbErr := database.DB.Exec(
				"INSERT INTO checks (website_id, status, response_time_ms) VALUES ($1, $2, $3)",
				targetWeb.ID, status, latency,
			)
			if dbErr != nil {
				slog.Error("Gagal mencatat log", "url", targetWeb.URL, "error", dbErr, "source", "manual_check")
			}

			resultsChan <- models.WebsiteStatus{URL: targetWeb.URL, Status: status, ResponseTimeMs: latency}
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

	c.JSON(http.StatusOK, checkResults)
}
