package models

// Define a GORM model
type Game struct {
	ID         string
	Difficulty int
	Status     string
	Players    []Player
}
