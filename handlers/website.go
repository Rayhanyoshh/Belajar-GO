package handlers

import (
	"belajar-go/database"
	"belajar-go/models"
	"net/http"
	"time"

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

// DeleteWebsite godoc
// @Summary      Menghapus website
// @Description  Menghapus URL dari daftar pantauan berdasarkan ID.
// @Tags         Websites
// @Produce      json
// @Param        id   path      int  true  "Website ID"
// @Success      200  {object}  map[string]string
// @Router       /websites/{id} [delete]
// @Security     BearerAuth
func DeleteWebsite(c *gin.Context) {
	id := c.Param("id")

	// Pastikan data ada
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM websites WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Website tidak ditemukan"})
		return
	}

	// Hapus checks terkait terlebih dahulu (karena ada foreign key constraint)
	_, err = database.DB.Exec("DELETE FROM checks WHERE website_id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus riwayat check website"})
		return
	}

	// Hapus website
	_, err = database.DB.Exec("DELETE FROM websites WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus website"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Website berhasil dihapus"})
}

// UpdateWebsite godoc
// @Summary      Memperbarui website
// @Description  Mengubah URL website yang sedang dipantau berdasarkan ID.
// @Tags         Websites
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Website ID"
// @Param        request body AddWebsiteInput true "Data URL baru"
// @Success      200  {object}  models.Website
// @Router       /websites/{id} [put]
// @Security     BearerAuth
func UpdateWebsite(c *gin.Context) {
	id := c.Param("id")
	var input AddWebsiteInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON salah"})
		return
	}

	// Pastikan data ada
	var web models.Website
	err := database.DB.QueryRow("SELECT id FROM websites WHERE id = $1", id).Scan(&web.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Website tidak ditemukan"})
		return
	}

	// Update website
	_, err = database.DB.Exec("UPDATE websites SET url = $1 WHERE id = $2", input.URL, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui website (mungkin URL sudah digunakan)"})
		return
	}

	web.URL = input.URL
	c.JSON(http.StatusOK, web)
}

// GetWebsiteHistory godoc
// @Summary      Riwayat check website
// @Description  Mendapatkan 50 riwayat pengecekan terakhir untuk suatu website.
// @Tags         Websites
// @Produce      json
// @Param        id   path      int  true  "Website ID"
// @Success      200  {array}   models.CheckHistory
// @Router       /websites/{id}/history [get]
// @Security     BearerAuth
func GetWebsiteHistory(c *gin.Context) {
	id := c.Param("id")

	// Pastikan website ada
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM websites WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Website tidak ditemukan"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT checked_at, status, response_time_ms 
		FROM checks 
		WHERE website_id = $1 
		ORDER BY checked_at DESC 
		LIMIT 50
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat dari database"})
		return
	}
	defer rows.Close()

	var history []models.CheckHistory
	for rows.Next() {
		var h models.CheckHistory
		var checkedAt time.Time
		
		if err := rows.Scan(&checkedAt, &h.Status, &h.ResponseTimeMs); err != nil {
			continue
		}
		
		h.CheckedAt = checkedAt.Format(time.RFC3339)
		history = append(history, h)
	}

	if history == nil {
		history = []models.CheckHistory{}
	}

	c.JSON(http.StatusOK, history)
}

