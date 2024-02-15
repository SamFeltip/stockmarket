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

func (game Game) GenerateInsights(db *gorm.DB) error {
	// get every game stock, and pull the insights for each stock
	// then distribute them randomly to players

	// delete all existing insights for current game (player insights where the player.gameID is this game)
	// get every player insight for this game

	// get all player stocks where playerStock.gameStock.gameID is this game
	// based on this sql:
	/*
		SELECT player_stocks.*
		FROM player_stocks
		INNER JOIN game_stocks ON player_stocks.game_stock_id = game_stocks.id
		WHERE game_stocks.game_id = 'your_game_id';
	*/
	var playerStocks []PlayerStock
	err := db.Where("game_stocks.game_id = ?", game.ID).Find(&playerStocks).Error

	if err != nil {
		return err
	}

	// get all insights
	var insights []Insight
	err = db.Find(&insights).Error

	if err != nil {
		return err
	}

	// for each game stock, distribute insights to players
	return nil
}
