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
}

type PlayerStockDisplay struct {
	ID    uint
	Value float64 //sum of game stock value and player stock quantity
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

	var playerStocksResult = []PlayerStockDisplay{}

	err := db.Table("player_stocks as ps").
		Select("ps.id, (ps.quantity * gs.value) as value").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Joins("inner join stocks as s on s.id = gs.stock_id").
		Where("ps.player_id = ?", playerID).
		Order("s.variation").
		Scan(&playerStocksResult).Error

	if err != nil {
		fmt.Println("could not load player stocks", err)
		return CurrentPlayerDisplay{}, err
	}

	if len(playerStocksResult) == 0 {
		fmt.Println("no player stocks found for this player")
		return CurrentPlayerDisplay{}, gorm.ErrRecordNotFound
	}

	var currentPlayerResult = struct {
		ID              uint
		UserID          uint
		UserName        string
		UserProfileRoot string
		Cash            int
		Active          bool
	}{}

	err = db.Table("players as p").
		Select("p.id, p.user_id, u.name as user_name, u.profile_root as user_profile_root, p.cash, p.active").
		Joins("inner join users as u on p.user_id = u.id").
		Where("p.id = ?", playerID).
		First(&currentPlayerResult).Error

	var currentPlayer = CurrentPlayerDisplay{
		ID:              currentPlayerResult.ID,
		UserID:          currentPlayerResult.UserID,
		UserName:        currentPlayerResult.UserName,
		UserProfileRoot: currentPlayerResult.UserProfileRoot,
		Cash:            currentPlayerResult.Cash,
		Active:          currentPlayerResult.Active,
		PlayerStocks:    playerStocksResult,
	}

	return currentPlayer, err

	/*
		var currentPlayerResult = []struct {
			UserID           uint
			UserName         string
			UserProfileRoot  string
			Cash             int
			Active           bool
			PlayerID         uint
			PlayerStockID    uint
			PlayerStockValue float64
		}{}

		err := db.Table("players as p").
			Select(
				"p.user_id, u.name as user_name, "+
					"u.profile_root as user_profile_root, "+
					"p.cash, ps.id as player_stock_id, "+
					"p.id as player_id, ps.id as player_stock_id,"+
					"(ps.quantity * gs.value) as player_stock_value").
			Joins("inner join users as u on p.user_id = u.id").
			Joins("inner join player_stocks as ps on ps.player_id = p.id").
			Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
			Where("p.id = ?", playerID).
			Scan(&currentPlayerResult).Error

		if err != nil {
			fmt.Println("could not load current player display", err)
			return CurrentPlayerDisplay{}, err
		}

		var playerDisplay = CurrentPlayerDisplay{
			ID:              currentPlayerResult[0].PlayerID,
			UserID:          currentPlayerResult[0].UserID,
			UserName:        currentPlayerResult[0].UserName,
			UserProfileRoot: currentPlayerResult[0].UserProfileRoot,
			Cash:            currentPlayerResult[0].Cash,
			Active:          currentPlayerResult[0].Active,
		}

		if len(playerDisplay.PlayerStocks) > 0 {
			fmt.Println("could not load player stocks")
			return playerDisplay, gorm.ErrRecordNotFound
		}

		return playerDisplay, err

	*/
}

func (playerDisplay *CurrentPlayerDisplay) TotalValue() float64 {

	var total float64
	for _, ps := range playerDisplay.PlayerStocks {
		total += float64(ps.Value)
	}
	return total + float64(playerDisplay.Cash)
}
