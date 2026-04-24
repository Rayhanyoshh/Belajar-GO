package handlers

import (
	"belajar-go/database"
	"belajar-go/models"
	"fmt"
	"net/http"
	"sync"

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

	c.JSON(http.StatusOK, checkResults)
}
