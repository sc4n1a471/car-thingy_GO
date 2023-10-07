package models

type InspectionResult struct {
	LicensePlate string   `json:"license_plate"`
	Name         string   `json:"name"`
	Base64       []string `json:"base_64"`
}
