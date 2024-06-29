package models

type Inspection struct {
	LicensePlate  string `json:"licensePlate"`
	Name          string `json:"name"`
	ImageLocation string `json:"imageLocation"`
}
