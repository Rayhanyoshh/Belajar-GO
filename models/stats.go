package models

type Stats struct {
	TotalWebsites    int     `json:"total_websites"`
	WebsitesUp       int     `json:"websites_up"`
	WebsitesDown     int     `json:"websites_down"`
	UptimePercentage float64 `json:"uptime_percentage"`
	TotalChecksToday int     `json:"total_checks_today"`
}
