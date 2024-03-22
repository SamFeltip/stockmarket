package models

import "gorm.io/gorm"

type PlayerDisplay struct {
	PlayerID     uint
	UserID       uint
	Name         string
	ProfileRoot  string
	Cash         int
	PlayerStocks []PlayerStockDisplay `gorm:"foreignKey:PlayerID"`
}

type PlayerStockDisplay struct {
	ID    uint
	Value float64 //sum of game stock value and player stock quantity
}

func LoadPlayerDisplays(gameID string, db *gorm.DB) ([]PlayerDisplay, error) {

	var players []PlayerDisplay
	err := db.Table("players").
		Select("u.name as user_name, u.profile_root as user_profile_root, p.cash").
		Joins("inner join users as u on u.id = players.user_id").
		Where("game_id = ?", gameID).
		Scan(&players).
		Error

	return players, err
}

func LoadCurrentPlayerDisplay(playerID uint, db *gorm.DB) (PlayerDisplay, error) {

	var playerDisplay PlayerDisplay

	err := db.Table("players as p").
		Select("p.id as player_id, p.user_id, u.name, u.profile_root, p.cash, ps.id as player_stock_id, (ps.quantity * gs.value) as player_stock_value").
		Joins("inner join users as u on p.user_id = u.id").
		Joins("inner join player_stocks as ps on ps.player_id = p.id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Where("p.id = ?", playerID).
		First(&playerDisplay).Error

	return playerDisplay, err
}

func (playerDisplay *PlayerDisplay) TotalValue() float64 {

	var total float64
	for _, ps := range playerDisplay.PlayerStocks {
		total += float64(ps.Value)
	}
	return total + float64(playerDisplay.Cash)
}
