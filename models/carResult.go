package models

type CarResult struct {
	LicensePlate LicensePlate       `json:"licensePlate"`
	Accidents    []Accident         `json:"accidents,omitempty"`
	Specs        Specs              `json:"specs"`
	Restrictions []Restriction      `json:"restrictions,omitempty"`
	Mileage      []Mileage          `json:"mileage"`
	Coordinates  Coordinate         `json:"coordinates"`
	Inspections  []InspectionResult `json:"inspections,omitempty"`
}
