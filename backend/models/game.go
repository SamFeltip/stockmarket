package models

// Define a GORM model
type Game struct {
	ID         string `gorm:"primaryKey"`
	Difficulty int
	Status     string
	Players    []Player
}
