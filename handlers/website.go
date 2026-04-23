package handlers

import (
	"belajar-go/database"
	"belajar-go/models"
	"encoding/json"
	"net/http"
)

// GetWebsites godoc
// @Summary      Mengambil daftar website
// @Description  Mendapatkan semua URL website yang sedang dipantau oleh sistem.
// @Tags         Websites
// @Produce      json
// @Success      200  {array}   models.Website
// @Router       /websites [get]
func GetWebsites(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, url FROM websites")
	if err != nil {
		http.Error(w, "Gagal mengambil data dari database", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(websites)
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
func AddWebsite(w http.ResponseWriter, r *http.Request) {
	var input AddWebsiteInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Format JSON salah", http.StatusBadRequest)
		return
	}

	var newID int
	err := database.DB.QueryRow(
		"INSERT INTO websites (url) VALUES ($1) RETURNING id", 
		input.URL,
	).Scan(&newID)

	if err != nil {
		http.Error(w, "Gagal menyimpan ke database (mungkin URL sudah ada)", http.StatusInternalServerError)
		return
	}

	newWeb := models.Website{
		ID:  newID,
		URL: input.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newWeb)
}
