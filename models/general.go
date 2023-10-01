package models

import "time"

type General struct {
	LicensePlate string    `json:"license_plate"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
}
