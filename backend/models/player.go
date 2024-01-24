package models

import "gorm.io/gorm"

type Player struct {
	ID       uint
	GameID   string
	Game     Game
	UserID   uint
	User     User
	Position int
}

func GetPlayer(game Game, user User, db *gorm.DB) (Player, error) {

	var player Player
	err := db.Where("game_id = ? AND user_id = ?", game.ID, user.ID).First(&player).Error

	return player, err
}