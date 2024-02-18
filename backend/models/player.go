package models

import (
	"fmt"
	"sort"

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

func GetPlayer(game *Game, user *User, db *gorm.DB) (Player, error) {

	var player Player
	err := db.Where("game_id = ? AND user_id = ?", game.ID, user.ID).First(&player).Error

	return player, err
}

func PlayerLeft(userID uint, gameID string, db *gorm.DB) error {

	var player Player
	err := db.Where("game_id = ? AND user_id = ?", gameID, userID).First(&player).Error

	if err != nil {
		fmt.Println("could not fetch player")
		return err
	}

	player.Active = false
	err = db.Save(&player).Error

	if err != nil {
		fmt.Println("could not update player to inactive")
		return err
	}

	return err
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

func (player *Player) TotalValue() float64 {

	var total float64
	for _, ps := range player.PlayerStocks {
		total += float64(ps.Quantity) * ps.GameStock.Value
	}
	return total + float64(player.Cash)
}

func (player Player) SortedPlayerStocks() []PlayerStock {

	player_stocks := player.PlayerStocks

	sort.Slice(player_stocks, func(i, j int) bool {
		return player_stocks[i].GameStock.Stock.Variation < player_stocks[j].GameStock.Stock.Variation
	})

	return player_stocks
}
