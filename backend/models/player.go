package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	ID       uint
	GameID   string
	Game     Game
	UserID   uint
	User     User
	Position int
	Active   bool
}

func GetPlayer(game *Game, user *User, db *gorm.DB) (Player, error) {

	var player Player
	err := db.Where("game_id = ? AND user_id = ?", game.ID, user.ID).First(&player).Error

	return player, err
}

func PlayerLeft(userID uint, gameID string, db *gorm.DB) (Player, error) {

	var player Player
	err := db.Where("game_id = ? AND user_id = ?", gameID, userID).First(&player).Error

	if err != nil {
		fmt.Println("could not fetch player")
		return Player{}, err
	}

	player.Active = false
	err = db.Save(&player).Error

	if err != nil {
		fmt.Println("could not update player to inactive")
		return Player{}, err
	}

	return player, err
}

// sort array of players by ID bubble sort
func SortPlayers(players []Player) []Player {

	for i := 0; i < len(players); i++ {
		for j := 0; j < len(players)-i-1; j++ {
			if players[j].ID > players[j+1].ID {
				players[j], players[j+1] = players[j+1], players[j]
			}
		}
	}

	return players
}
