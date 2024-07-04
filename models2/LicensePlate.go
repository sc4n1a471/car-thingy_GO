package models

type LicensePlate struct {
	LicensePlate string `json:"licensePlate" gorm:"primaryKey"`
	Comment      string `json:"comment"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
