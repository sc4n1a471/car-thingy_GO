package models

type General struct {
	LicensePlate string  `json:"license_plate"`
	latitude     float32 `json:"latitude"`
	longitude    float32 `json:"longitude"`
	Comment      string  `json:"Comment"`
}
