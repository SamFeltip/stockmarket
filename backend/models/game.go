package models

import (
	"errors"
	"fmt"
	"math/rand"

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

func (game Game) GetPlayer(user User) (Player, error) {

	for _, player := range game.Players {
		if player.UserID == user.ID {
			return player, nil
		}
	}

	return Player{}, errors.New("Player not found")
}

func (game Game) MustGetPlayer(user User) Player {
	player, err := game.GetPlayer(user)
	if err != nil {
		fmt.Println("could not must get player, return nil")
		return Player{}
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

	fmt.Println("game.ID:", game.ID)

	player_insights := []PlayerInsight{}
	err := db.Joins("INNER JOIN player_stocks on player_stocks.id = player_insights.player_stock_id").
		Joins("INNER JOIN game_stocks on game_stocks.id = player_stocks.game_stock_id").
		Where("game_stocks.game_id = ?", game.ID).Find(&player_insights).Error

	if err != nil {
		fmt.Println("could not get player insights", err)
		return err
	}

	fmt.Println("player_insights.len:", len(player_insights))

	err = db.Delete(&player_insights).Error

	if err != nil {
		fmt.Println("could not find  player insights to delete", err)
	} else {
		fmt.Println("deleted player insights")
	}

	// get all insights
	var insights []Insight
	err = db.Preload("Stock").Find(&insights).Error

	if err != nil {
		fmt.Println("could not get insights", err)
		return err
	}

	fmt.Println("insights.len:", len(insights))

	// get num_of_players * 10 random insights
	var players []Player
	db.Where("game_id = ?", game.ID).Find(&players)

	player_length := len(players)

	// shuffle insights
	rand.Shuffle(len(insights), func(i, j int) {
		insights[i], insights[j] = insights[j], insights[i]
	})

	fmt.Println("shuffled insights")

	top_insights := insights[:player_length*10]

	// give 10 playerInsights to each player for each insight
	starting_point := 0
	for _, player := range players {

		for i := starting_point; i < starting_point+10; i++ {
			// get player_stock for player of top_insights[i].Stock
			player_stock := PlayerStock{}
			err = db.Joins("INNER JOIN game_stocks on player_stocks.game_stock_id = game_stocks.id").
				Where("game_stocks.stock_id = ? AND player_id = ?", top_insights[i].Stock.ID, player.ID).First(&player_stock).Error

			if err != nil {
				fmt.Println("could not get player stock", err)
				continue
			}

			fmt.Println("got player stock")

			// create player_insight
			player_insight := PlayerInsight{
				PlayerStock: player_stock,
				Insight:     top_insights[i],
			}

			err = db.Create(&player_insight).Error

			if err != nil {
				fmt.Println("could not create player insight", err)
			}
		}

		starting_point += 10
	}

	return nil
}
