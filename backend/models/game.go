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

type GameStatus string

var Waiting GameStatus = "waiting"
var Playing GameStatus = "playing"
var Evaluating GameStatus = "evaluating"
var Finished GameStatus = "finished"

func GetGame(gameID string, db *gorm.DB) (Game, error) {

	var game Game
	err := db.Model(&game).Preload("CurrentUser").Preload("Players").Preload("Players.User").Where("lower(id) = lower(?)", gameID).First(&game).Error

	return game, err

}

func (game Game) UpdateORM(db *gorm.DB) error {
	err := db.Model(&game).Preload("CurrentUser").Preload("Players").Preload("Players.User").Where("lower(id) = lower(?)", game.ID).First(&game).Error
	return err
}

func GameDifficultyDisplay(difficulty int) string {
	switch difficulty {
	case 0:
		return "Short"
	case 1:
		return "Medium"
	case 2:
		return "Long"
	default:
		return "Unknown"
	}
}
