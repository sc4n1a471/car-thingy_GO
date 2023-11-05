package models

type LicensePlate struct {
	LicensePlate string `json:"license_plate" gorm:"primaryKey"`
	Comment      string `json:"comment"`
	CreatedAt    string `json:"created_at"`
}
