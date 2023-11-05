package models

type Accident struct {
	LicensePlate string `json:"license_plate"`
	AccidentDate string `json:"accident_date"`
	Role         string `json:"role"`
}
