package models

type Mileage struct {
	LicensePlate string `json:"licensePlate"`
	Mileage      int    `json:"mileage"`
	MileageDate  string `json:"mileageDate"`
}
