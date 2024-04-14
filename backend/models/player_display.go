package models

import (
	"fmt"

	"gorm.io/gorm"
)

type PlayerDisplay struct {
	PlayerID        uint
	UserID          uint
	UserName        string
	UserProfileRoot string
	Cash            int
	Active          bool
	TotalValue      float64
	Ranking         int
}

type CurrentPlayerDisplay struct {
	ID              uint
	UserID          uint
	UserName        string
	UserProfileRoot string
	Cash            int
	Active          bool
	PlayerStocks    []PlayerStockDisplay
	TotalValue      float64
}

func LoadPlayerDisplays(gameID string, db *gorm.DB) ([]PlayerDisplay, error) {

	var players []PlayerDisplay
	err := db.Table("players").
		Select("players.ID as player_id, u.ID as user_id, u.name as user_name, u.profile_root as user_profile_root, cash, active").
		Joins("inner join users as u on u.id = players.user_id").
		Where("game_id = ?", gameID).
		Scan(&players).
		Error

	return players, err
}

func LoadCurrentPlayerDisplay(playerID uint, db *gorm.DB) (CurrentPlayerDisplay, error) {

	var currentPlayerResult struct {
		GameID          string
		UserID          uint
		UserName        string
		UserProfileRoot string
		Active          bool
		Cash            int
		NetValue        float64
	}

	err := db.Table("players as p").
		Select("p.game_id, p.user_id, u.name as user_name, u.profile_root as user_profile_root, p.active, p.cash, sum(ps.quantity * gs.value) as net_value").
		Joins("inner join users as u on p.user_id = u.id").
		Joins("inner join player_stocks as ps on ps.player_id = p.id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Where("p.id = ?", playerID).
		Group("p.game_id, p.user_id, u.name, u.profile_root, p.cash, p.active").
		First(&currentPlayerResult).Error

	if err != nil {
		fmt.Println("could not load current player display", err)
		return CurrentPlayerDisplay{}, err
	}

	playerStockDisplays, err := GetPlayerStockPreviews(playerID, db)

	if err != nil {
		fmt.Println("could not load player stocks", err)
		return CurrentPlayerDisplay{}, err
	}

	var currentPlayer = CurrentPlayerDisplay{
		ID:              playerID,
		UserID:          currentPlayerResult.UserID,
		UserName:        currentPlayerResult.UserName,
		UserProfileRoot: currentPlayerResult.UserProfileRoot,
		Cash:            currentPlayerResult.Cash,
		TotalValue:      currentPlayerResult.NetValue,
		Active:          currentPlayerResult.Active,
		PlayerStocks:    playerStockDisplays,
	}

	return currentPlayer, nil
}

func LoadPlayerDisplay(id uint, db *gorm.DB) (PlayerDisplay, error) {

	var currentPlayerResult struct {
		GameID          string
		UserID          uint
		UserName        string
		UserProfileRoot string
		Active          bool
		Cash            int
		NetValue        float64
	}

	err := db.Table("players as p").
		Select("p.game_id, p.user_id, u.name as user_name, u.profile_root as user_profile_root, p.active, p.cash, sum(ps.quantity * gs.value) as net_value").
		Joins("inner join users as u on p.user_id = u.id").
		Joins("inner join player_stocks as ps on ps.player_id = p.id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Where("p.id = ?", id).
		Group("p.game_id, p.user_id, u.name, u.profile_root, p.cash, p.active").
		First(&currentPlayerResult).Error

	if err != nil {
		fmt.Println("could not load current player display", err)
		return PlayerDisplay{}, err
	}

	type PlayerPositionResponse struct {
		ID         uint
		TotalValue float64
	}

	player_positions := []PlayerPositionResponse{}

	error := db.Table("players as p").
		Select("p.id, (p.cash + sum(ps.quantity * gs.value)) as total_value").
		Joins("inner join player_stocks as ps on ps.player_id = p.id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Where("p.game_id = ?", currentPlayerResult.GameID).
		Group("p.id, p.game_id, p.user_id, p.cash").
		Order("total_value").
		Scan(&player_positions).Error

	if error != nil {
		fmt.Println("could not load player positions", error)
		return PlayerDisplay{}, error
	}

	ranking := 0

	for i, player := range player_positions {
		if player.ID == id {
			ranking = len(player_positions) - i
		}
	}

	return PlayerDisplay{
		PlayerID:        id,
		UserID:          currentPlayerResult.UserID,
		UserName:        currentPlayerResult.UserName,
		UserProfileRoot: currentPlayerResult.UserProfileRoot,
		Cash:            currentPlayerResult.Cash,
		TotalValue:      currentPlayerResult.NetValue,
		Ranking:         ranking,
	}, nil
}

func (playerDisplay *PlayerDisplay) TotalShares(playerStocks []PlayerStockDisplay) int {
	totalShares := 0
	for _, playerStock := range playerStocks {
		totalShares += playerStock.PlayerStockQuantity
	}
	return totalShares
}
