package models

type Restriction struct {
	LicensePlate    string `json:"license_plate"`
	Restriction     string `json:"restriction"`
	RestrictionDate string `json:"restriction_date"`
	Active          bool   `json:"active"`
}
