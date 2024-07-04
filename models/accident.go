package models

type Accident struct {
	ID           int    `json:"id" gorm:"primaryKey,autoIncrement"`
	CarID        string `json:"licensePlate" gorm:"size:255"`
	AccidentDate string `json:"accidentDate"`
	Role         string `json:"role"`
}
