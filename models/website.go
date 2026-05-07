package models

// Ini adalah cetakan data yang akan bolak-balik antara JSON dan Database

type Website struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type WebsiteStatus struct {
	URL            string `json:"url"`
	Status         string `json:"status"` // "UP" atau "DOWN"
	ResponseTimeMs int64  `json:"response_time_ms"`
}

type CheckHistory struct {
	CheckedAt      string `json:"checked_at"`
	Status         string `json:"status"`
	ResponseTimeMs *int64 `json:"response_time_ms"` // Pointer agar bisa handle null dari data lama
}
