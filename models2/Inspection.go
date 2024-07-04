package models

type Inspection struct {
	LicensePlate  string `json:"licensePlate"`
	Name          string `json:"name"`
	ImageLocation string `json:"imageLocation"`
}

type InspectionResult struct {
	LicensePlate string   `json:"licensePlate"`
	Name         string   `json:"name"`
	Base64       []string `json:"base64"`
}
