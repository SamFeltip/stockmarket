package models

// Define a GORM model
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
}
