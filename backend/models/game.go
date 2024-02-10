package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Define a GORM model
type Game struct {
	gorm.Model
	ID            string
	Difficulty    int
	Status        string
	Players       []Player
	GameStocks    []GameStock
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
	err := db.Model(&game).
		Preload("GameStocks").
		Preload("GameStocks.Stock").
		Preload("CurrentUser").
		Preload("Players").
		Preload("Players.User").
		Preload("Players.PlayerStocks").
		Preload("Players.PlayerStocks.GameStock.Stock").
		Where("lower(id) = lower(?)", gameID).First(&game).Error

	return game, err

}

func (game Game) UpdateORM(db *gorm.DB) error {
	err := db.Model(&game).Preload("CurrentUser").Preload("Players").Preload("Players.User").Where("lower(id) = lower(?)", game.ID).First(&game).Error
	return err
}

func (game Game) GetPlayer(user User) (*Player, error) {

	for _, player := range game.Players {
		if player.UserID == user.ID {
			return &player, nil
		}
	}

	return nil, errors.New("Player not found")
}

func (game Game) MustGetPlayer(user User) *Player {
	player, err := game.GetPlayer(user)
	if err != nil {
		fmt.Println("could not must get player, return nil")
		return nil
	}
	return player
}

func (game Game) UpdateStatus(status GameStatus, db *gorm.DB) error {
	err := db.Model(&game).Where("id = lower(?)", game.ID).Update("status", status).Error
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
