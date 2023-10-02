package models

import "time"

type Specs struct {
	LicensePlate string    `json:"license_plate" gorm:"primaryKey"`
	Brand        string    `json:"brand"`
	Color        string    `json:"color"`
	EngineSize   int       `json:"engine_size"`
	FirstReg     time.Time `json:"first_reg"`
	FirstRegHun  time.Time `json:"first_reg_hun"`
	FuelType     string    `json:"fuel_type"`
	Gearbox      string    `json:"gearbox"`
	Model        string    `json:"model"`
	NumOfOwners  int       `json:"num_of_owners"`
	Performance  int       `json:"performance"`
	Status       string    `json:"status"`
	TypeCode     string    `json:"type_code"`
	Year         int       `json:"year"`
	Comment      string    `json:"comment"`
}
