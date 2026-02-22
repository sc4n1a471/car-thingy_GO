package models

type QueryLog struct {
	User           string `json:"user"`
	LicensePlate   string `json:"licensePlate"`
	QueryTimestamp string `json:"queryTimestamp"`
}
