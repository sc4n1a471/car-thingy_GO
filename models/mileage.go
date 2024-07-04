package models

type Mileage struct {
	ID      int    `json:"id" gorm:"primaryKey,autoIncrement"`
	CarID   string `json:"licensePlate" gorm:"size:255"`
	Mileage int    `json:"mileage"`
	Date    string `json:"date"`
}
