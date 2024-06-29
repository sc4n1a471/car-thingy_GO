package models

type InspectionResult struct {
	LicensePlate string   `json:"licensePlate"`
	Name         string   `json:"name"`
	Base64       []string `json:"base64"`
}
