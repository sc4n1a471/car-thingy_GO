package models

type Inspection struct {
	ID            int      `json:"id" gorm:"primaryKey;autoIncrement"`
	CarID         string   `json:"licensePlate" gorm:"size:255"`
	ImageLocation string   `json:"imageLocation"`
	Name          string   `json:"name"`
	Base64        []string `json:"base64" gorm:"-"`
}

type QueryInspection struct {
	ID            int      `json:"id" gorm:"primaryKey,autoIncrement`
	CarID         string   `json:"licensePlate"`
	ImageLocation string   `json:"imageLocation"`
	Name          string   `json:"name"`
	Base64        []string `json:"base64" gorm:"-"`
}
