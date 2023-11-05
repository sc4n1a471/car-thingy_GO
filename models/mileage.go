package models

type Mileage struct {
	LicensePlate string `json:"license_plate"`
	Mileage      int    `json:"mileage"`
	MileageDate  string `json:"mileage_date"`
}
