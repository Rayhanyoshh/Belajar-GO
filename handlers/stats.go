package handlers

import (
	"belajar-go/database"
	"belajar-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStats godoc
// @Summary      Statistik Dashboard
// @Description  Mendapatkan ringkasan statistik (jumlah website, uptime, dsb) untuk dashboard.
// @Tags         Stats
// @Produce      json
// @Success      200  {object}  models.Stats
// @Router       /stats [get]
func GetStats(c *gin.Context) {
	var stats models.Stats

	// 1. Total websites
	database.DB.QueryRow("SELECT COUNT(*) FROM websites").Scan(&stats.TotalWebsites)

	// Jika belum ada website, return kosong
	if stats.TotalWebsites == 0 {
		c.JSON(http.StatusOK, stats)
		return
	}

	// 2. Websites UP & DOWN (Berdasarkan check terbaru per website)
	// Query ini mencari status terbaru untuk setiap website
	rows, err := database.DB.Query(`
		SELECT status 
		FROM (
			SELECT website_id, status,
				ROW_NUMBER() OVER(PARTITION BY website_id ORDER BY checked_at DESC) as rn
			FROM checks
		) latest_checks
		WHERE rn = 1
	`)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var status string
			if err := rows.Scan(&status); err == nil {
				if status == "UP" {
					stats.WebsitesUp++
				} else {
					stats.WebsitesDown++
				}
			}
		}
	}

	// 3. Kalkulasi Uptime Percentage
	if stats.TotalWebsites > 0 {
		stats.UptimePercentage = float64(stats.WebsitesUp) / float64(stats.TotalWebsites) * 100
	}

	// 4. Total Checks Today
	database.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM checks 
		WHERE DATE(checked_at) = CURRENT_DATE
	`).Scan(&stats.TotalChecksToday)

	c.JSON(http.StatusOK, stats)
}
