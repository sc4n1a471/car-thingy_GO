package models

type Restriction struct {
	ID          int    `json:"id" gorm:"primaryKey,autoIncrement"`
	CarID       string `json:"licensePlate" gorm:"size:255"`
	Restriction string `json:"restriction"`
	IsActive    bool   `json:"isActive" gorm:"default:true"`
}
