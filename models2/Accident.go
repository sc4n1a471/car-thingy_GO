package models

type Accident struct {
	LicensePlate string `json:"licensePlate"`
	AccidentDate string `json:"accidentDate"`
	Role         string `json:"role"`
}
