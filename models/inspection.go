package models

type Inspection struct {
	LicensePlate  string `json:"license_plate"`
	Name          string `json:"name"`
	ImageLocation string `json:"image_location"`
}
