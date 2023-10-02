package models

type CarResult struct {
	Accidents    []Accident         `json:"accidents"`
	Specs        Specs              `json:"specs"`
	Restrictions []Restriction      `json:"restrictions"`
	Mileage      []Mileage          `json:"mileage"`
	General      General            `json:"general"`
	Inspections  []InspectionResult `json:"inspections"`
}
