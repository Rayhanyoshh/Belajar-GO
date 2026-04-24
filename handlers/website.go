package handlers

import (
	"belajar-go/database"
	"belajar-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetWebsites godoc
// @Summary      Mengambil daftar website
// @Description  Mendapatkan semua URL website yang sedang dipantau oleh sistem.
// @Tags         Websites
// @Produce      json
// @Success      200  {array}   models.Website
// @Router       /websites [get]
func GetWebsites(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, url FROM websites")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dari database"})
		return
	}
	defer rows.Close()

	var websites []models.Website
	for rows.Next() {
		var web models.Website
		if err := rows.Scan(&web.ID, &web.URL); err != nil {
			continue
		}
		websites = append(websites, web)
	}

	if websites == nil {
		websites = []models.Website{}
	}

	c.JSON(http.StatusOK, websites)
}

// Struct tambahan khusus untuk memberi tahu dokumentasi wujud input JSON-nya
type AddWebsiteInput struct {
	URL string `json:"url" example:"https://google.com"`
}

// AddWebsite godoc
// @Summary      Mendaftarkan website baru
// @Description  Menambahkan URL baru ke dalam database untuk dipantau secara otomatis.
// @Tags         Websites
// @Accept       json
// @Produce      json
// @Param        request body AddWebsiteInput true "Data URL yang ingin dipantau"
// @Success      201  {object}  models.Website
// @Router       /websites [post]
func AddWebsite(c *gin.Context) {
	var input AddWebsiteInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON salah"})
		return
	}

	var newID int
	err := database.DB.QueryRow(
		"INSERT INTO websites (url) VALUES ($1) RETURNING id", 
		input.URL,
	).Scan(&newID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database (mungkin URL sudah ada)"})
		return
	}

	newWeb := models.Website{
		ID:  newID,
		URL: input.URL,
	}

	c.JSON(http.StatusCreated, newWeb)
}
