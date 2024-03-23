package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	ID           uint
	GameID       string
	Game         Game
	PlayerStocks []PlayerStock `gorm:"constraint:OnDelete:CASCADE"`
	UserID       uint
	User         User
	Cash         int
	Position     int
	Active       bool
}

func FindPlayer(playerID uint, db *gorm.DB) (Player, error) {

	var player Player
	err := db.First(&player, "id = ?", playerID).Error

	return player, err
}

func PlayerLeft(playerID uint, db *gorm.DB) error {

	var player Player
	err := db.Model(&player).Where("id = ?", playerID).Update("active", false).Error

	if err != nil {
		fmt.Println("could not update player to inactive")
		return err
	}

	return err
}
