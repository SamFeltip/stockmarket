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
	Players       []Player    `gorm:"constraint:OnDelete:CASCADE"`
	GameStocks    []GameStock `gorm:"constraint:OnDelete:CASCADE"`
	CurrentUser   User
	CurrentUserID uint
}

type GameStatus string

var Waiting GameStatus = "waiting"
var Playing GameStatus = "playing"
var Evaluating GameStatus = "evaluating"
var Finished GameStatus = "finished"

func LoadGame(gameID string, db *gorm.DB) (Game, error) {

	var game Game
	err := db.Model(&game).
		Preload("GameStocks").
		Preload("GameStocks.Stock").
		Preload("CurrentUser").
		Preload("Players").
		Preload("Players.User").
		Where("lower(games.id) = lower(?)", gameID).
		First(&game).Error

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

func (game *Game) UpdateCurrentUser(db *gorm.DB) error {

	current_user := game.CurrentUser
	players := game.Players

	fmt.Println("sorting players by id")
	// sort players by id
	for i := 0; i < len(players); i++ {
		for j := i + 1; j < len(players); j++ {
			if players[i].ID > players[j].ID {
				players[i], players[j] = players[j], players[i]
			}
		}
	}

	fmt.Println("finding the next user")
	next_user := User{}
	// find the next player in the list (based on current player)
	for i, player := range players {
		if player.User.ID == current_user.ID {
			if i == len(players)-1 {
				next_user = players[0].User
			} else {
				next_user = players[0].User
			}
			break
		}
	}

	fmt.Println("setting next user:", next_user.ID)
	err := db.Model(&game).Where("id = lower(?)", game.ID).Update("current_user_id", next_user.ID).Error

	if err != nil {
		fmt.Println("could not update current player", err)
		return err
	}

	game.CurrentUser = next_user
	game.CurrentUserID = next_user.ID

	return nil
}
