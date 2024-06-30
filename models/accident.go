package models

import (
	"time"
)

type Accident struct {
	ID           int       `json:"id" gorm:"primaryKey,autoIncrement"`
	CarID        string    `json:"licensePlate" gorm:"size:255"`
	AccidentDate time.Time `json:"accidentDate"`
	Role         string    `json:"role"`
}
