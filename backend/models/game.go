package models

// Define a GORM model
type Game struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Difficulty string
}
