package models

import (
	"gorm.io/gorm"
)

// Define a GORM model
type Game struct {
	gorm.Model
	ID            string
	Difficulty    int
	Status        string
	Players       []Player
	CurrentUser   User
	CurrentUserID uint
}

func GetGame(gameID string, db *gorm.DB) (Game, error) {

	var game Game
	err := db.Model(&game).Preload("Players").Preload("Players.User").Where("lower(id) = lower(?)", gameID).First(&game).Error

	return game, err

}
