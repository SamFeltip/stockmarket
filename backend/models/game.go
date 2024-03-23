package models

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"gorm.io/gorm"
)

// Define a GORM model
type Game struct {
	gorm.Model
	ID            string
	PeriodCount   int
	CurrentPeriod int
	Status        string
	Players       []Player    `gorm:"constraint:OnDelete:CASCADE"`
	GameStocks    []GameStock `gorm:"constraint:OnDelete:CASCADE"`
	CurrentUser   User
	CurrentUserID uint
	FeedItems     []FeedItem `gorm:"constraint:OnDelete:CASCADE"`
}

type GameStatus string

var Waiting GameStatus = "waiting"
var Playing GameStatus = "playing"
var Closed GameStatus = "closed"
var Finished GameStatus = "finished"

/*
creates a new game and all possible game stocks

output: game (code from c post PeriodCount from gin context, user from gin context)

also runs CreateGameStocks
*/
func CreateGame(code string, periodCount int, current_user User, db *gorm.DB) (Game, error) {

	fmt.Println("create game:", code, periodCount)

	game := Game{
		ID:          code,
		PeriodCount: periodCount,
		Status:      string(Waiting),
		CurrentUser: current_user,
	}

	db.Create(&game)

	// get all stocks
	fmt.Println("create game stocks: ", code)
	game_stocks, err := CreateGameStocks(code, db)

	if err != nil {
		fmt.Println("error creating game stocks:", err)
		return Game{}, err
	}

	game.GameStocks = game_stocks
	db.Save(&game)

	if err != nil {
		fmt.Println("error setting active game:", err)
		return Game{}, err
	}

	return game, nil // passed into templates
}

func FindGame(gameID string, db *gorm.DB) (Game, error) {

	var game Game
	err := db.Where("lower(id) = lower(?)", gameID).First(&game).Error

	return game, err
}

func GamePeriodCountDisplay(periodCount int) string {
	switch periodCount {
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

/*
- set active game for passed in user

- if needed, create a player and associate it with the game
*/
func (current_user *User) SetActiveGame(gameID string, db *gorm.DB) (Player, error) {

	player, err := current_user.GetPlayer(gameID, db)

	// if gorm no record error
	if err != nil {
		fmt.Println("player does not exist, creating")
		player, err = current_user.CreatePlayer(gameID, db)

		if err != nil {
			fmt.Println("error creating player:", err)
			return Player{}, err
		}
	}

	if !player.Active {
		fmt.Println("setting active game for:", current_user.ID, gameID)

		err = db.Model(&player).Where("id = ?", player.ID).Update("active", true).Error
	}

	if err != nil {
		fmt.Println("error setting active game for:", current_user.ID, gameID, err)
		return Player{}, err
	}

	player.Active = true

	fmt.Println("unsetting active game for other games", current_user.ID, gameID)
	err = db.Model(&Player{}).Where("user_id = ? AND game_id != ?", current_user.ID, gameID).Update("active", false).Error

	if err != nil {
		fmt.Println("error unsetting active game for other games:", err)
		return Player{}, err
	}

	fmt.Println("set successfully")

	return player, nil
}

func (game *Game) GeneratePlayerInsights(players []Player, db *gorm.DB) error {

	fmt.Println("generating insights, game.ID:", game.ID)

	var playerInsightsIDs []int

	err := db.Table("player_insights").
		Select("player_insights.id").
		Joins("INNER JOIN player_stocks on player_stocks.id = player_insights.player_stock_id").
		Joins("INNER JOIN game_stocks on game_stocks.id = player_stocks.game_stock_id").
		Where("game_stocks.game_id = ?", game.ID).
		Model(&playerInsightsIDs).Error

	if err != nil {
		fmt.Println("could not get player insights", err)
		return err
	}

	player_insights := []PlayerInsight{}
	err = db.Where("id IN (?)", playerInsightsIDs).Delete(&player_insights).Error

	if err != nil {
		fmt.Println("could not find player insights to delete", err)
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

	// shuffle insights
	rand.Shuffle(len(insights), func(i, j int) {
		insights[i], insights[j] = insights[j], insights[i]
	})

	fmt.Println("shuffled insights")

	insights_per_player := 10

	if len(insights)/len(players) < insights_per_player {
		fmt.Println("not enough insights for all players")
		return errors.New("not enough insights for all players")
	}

	top_insights := insights[:len(players)*insights_per_player]

	// give insights_per_player playerInsights to each player for each insight
	starting_point := 0
	for _, player := range players {

		for i := starting_point; i < starting_point+insights_per_player; i++ {
			// get player_stock for player of top_insights[i].Stock
			player_stock := PlayerStock{}
			err = db.Joins("INNER JOIN game_stocks on player_stocks.game_stock_id = game_stocks.id").
				Where("game_stocks.stock_id = ? AND player_id = ?", top_insights[i].Stock.ID, player.ID).First(&player_stock).Error

			if err != nil {
				fmt.Println("could not get player stock", err)
				continue
			}

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

		starting_point += insights_per_player
	}

	return nil
}

/*
update the current user to the next user in the game
*/
func UpdateCurrentUser(gameID string, db *gorm.DB) (uint, error) {

	game, err := FindGame(gameID, db)

	if err != nil {
		fmt.Println("could not find game", err)
		return 0, err
	}

	var players = []struct {
		PlayerID uint
		UserID   uint
	}{}

	err = db.
		Table("players").
		Select("players.id as player_id, users.id as user_id").
		Joins("inner join users on players.user_id = users.id").
		Where("game_id = ? AND active = ?", gameID, true).
		Order("players.id").
		Find(&players).Error

	if err != nil {
		fmt.Println("could not get players", err)
		return 0, err
	}

	if len(players) == 0 {
		fmt.Println("no players in game")
		return 0, errors.New("no players in game")
	}

	current_user_id := game.CurrentUserID

	fmt.Println("finding the next user of ", strconv.Itoa(len(players)), " old user:", current_user_id)
	var next_user_id uint

	// find the next player in the list (based on current player)
	for i, player := range players {

		if player.UserID == current_user_id {
			if i == len(players)-1 {
				next_user_id = players[0].UserID
			} else {
				next_user_id = players[i+1].UserID
			}
			break
		}
	}

	game.CurrentUserID = next_user_id

	err = db.Save(&game).Error

	if err != nil {
		fmt.Println("could not update current player", err)
		return 0, err
	}

	return next_user_id, nil
}

func (game *Game) UpdatePeriod(db *gorm.DB) error {

	type GameStockChange struct {
		TotalChange float64
		GameStockID uint
		Value       float64
	}

	var gameStockChanges []GameStockChange

	err := db.Table("player_stocks as ps").
		Select("sum(i.value) as total_change, gs.id as game_stock_id, gs.value").
		Joins("left join player_insights as pi on pi.player_stock_id = ps.id").
		Joins("left join insights as i on i.id = pi.insight_id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Joins("inner join stocks as s on s.id = gs.stock_id").
		Where("gs.game_id = ?", game.ID).
		Group("gs.id, gs.value").
		Scan(&gameStockChanges).Error

	if err != nil {
		fmt.Println("could not get game stock changes", err)
		return err
	}

	// loop through gameStockChanges and update gameStocks
	for _, gameStockChange := range gameStockChanges {
		gameStock := GameStock{}
		err = db.
			Model(&gameStock).
			Where("id = ?", gameStockChange.GameStockID).
			Update("value", gameStockChange.Value+gameStockChange.TotalChange).Error

		if err != nil {
			fmt.Println("could not update game stock", err)
			return err
		}
	}

	players, err := GetPlayers(game.ID, db)

	if err != nil {
		fmt.Println("could not get players", err)
		return err
	}

	err = game.GeneratePlayerInsights(players, db)

	if err != nil {
		fmt.Println("could not generate player insights", err)
		return err
	}

	game.CurrentPeriod++
	game.Status = string(Playing)

	err = db.Save(&game).Error

	if err != nil {
		fmt.Println("could not update game", err)
		return err
	}

	// create feed item for new period
	feedItem := FeedItem{
		GameID: game.ID,
		Period: game.CurrentPeriod,

		Message:   "Players get another turn",
		Title:     "Period " + strconv.Itoa(game.CurrentPeriod+1),
		ImageRoot: "/static/imgs/icons/Stock.svg",
		Colour:    "grey",
	}

	err = db.Create(&feedItem).Error

	if err != nil {
		fmt.Println("could not create feed item", err)
		return err
	}

	return nil
}

func GetPlayers(gameID string, db *gorm.DB) ([]Player, error) {

	var players []Player
	err := db.Where("game_id = ? AND active = ?", gameID, true).Find(&players).Error

	return players, err
}
