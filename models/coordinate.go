package models

type Coordinate struct {
	LicensePlate string  `json:"license_plate" gorm:"primaryKey"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}
