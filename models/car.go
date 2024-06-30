package models

import (
	"time"
)

type Car struct {
	ID          string    `json:"licensePlate" gorm:"primaryKey"`
	Comment     *string   `json:"comment,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Brand       *string   `json:"brand,omitempty"`
	Color       *string   `json:"color,omitempty"`
	EngineSize  *int      `json:"engineSize,omitempty"`
	FirstReg    *string   `json:"firstReg,omitempty"`
	FirstRegHun *string   `json:"firstRegHun,omitempty"`
	FuelType    *string   `json:"fuelType,omitempty"`
	Gearbox     *string   `json:"gearbox,omitempty"`
	Model       *string   `json:"model,omitempty"`
	NumOfOwners *int      `json:"numOfOwners,omitempty"`
	Performance *int      `json:"performance,omitempty"`
	Status      *string   `json:"status,omitempty"`
	TypeCode    *string   `json:"typeCode,omitempty"`
	Year        *int      `json:"year,omitempty"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`

	Accidents    *[]Accident    `json:"accidents,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Inspections  *[]Inspection  `json:"inspections,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Mileage      *[]Mileage     `json:"mileage" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Restrictions *[]Restriction `json:"restrictions,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
