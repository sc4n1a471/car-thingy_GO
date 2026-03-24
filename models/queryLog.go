package models

import "time"

type QueryLog struct {
	User           string    `json:"user"`
	LicensePlate   string    `json:"licensePlate"`
	QueryTimestamp time.Time `json:"queryTimestamp"`
}
