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

type PlayerStockDisplay struct {
	ID             uint
	GameID         string
	TotalInsight   float64
	GameStockValue float64
	StockName      string
	StockImagePath string
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

	playerStockDisplays, err := GetPlayerStockDisplays(currentPlayerResult.GameID, playerID, db)

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
