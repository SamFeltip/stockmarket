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
