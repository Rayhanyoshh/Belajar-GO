package models

// Ini adalah cetakan data yang akan bolak-balik antara JSON dan Database

type Website struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type WebsiteStatus struct {
	URL    string `json:"url"`
	Status string `json:"status"` // "UP" atau "DOWN"
}
