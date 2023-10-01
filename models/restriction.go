package models

import "time"

type Restriction struct {
	LicensePlate    string    `json:"license_plate"`
	Restriction     string    `json:"restriction"`
	RestrictionDate time.Time `json:"restriction_date"`
	Active          bool      `json:"active"`
}
