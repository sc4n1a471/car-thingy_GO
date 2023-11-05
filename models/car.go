package models

type Car struct {
	LicensePlate LicensePlate  `json:"license_plate"`
	Accidents    []Accident    `json:"accidents"`
	Specs        Specs         `json:"specs"`
	Restrictions []Restriction `json:"restrictions"`
	Mileage      []Mileage     `json:"mileage"`
	Coordinates  Coordinate    `json:"coordinates"`
	Inspections  []Inspection  `json:"inspections"`
}
