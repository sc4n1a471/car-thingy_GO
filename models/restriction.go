package models

type Restriction struct {
	LicensePlate    string `json:"licensePlate"`
	Restriction     string `json:"restriction"`
	RestrictionDate string `json:"restrictionDate"`
	Active          bool   `json:"active"`
}
