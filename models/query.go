package models

type Car struct {
	Accidents    []Accident    `json:"accidents"`
	Specs        Specs         `json:"specs"`
	Restrictions []Restriction `json:"restrictions"`
	Mileage      []Mileage     `json:"mileage"`
}
