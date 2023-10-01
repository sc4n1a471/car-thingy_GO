package models

import "time"

type Accident struct {
	LicensePlate string    `json:"license_plate"`
	AccidentDate time.Time `json:"accident_date"`
	Role         string    `json:"role"`
}
