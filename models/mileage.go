package models

import "time"

type Mileage struct {
	LicensePlate string    `json:"license_plate"`
	Mileage      int       `json:"mileage"`
	MileageDate  time.Time `json:"mileage_date"`
}
