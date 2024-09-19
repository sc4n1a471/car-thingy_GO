package models

type AuthKey struct {
	ID       string `gorm:"primaryKey"`
	IsActive bool
}
