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

func LoadPlayer(playerID uint, db *gorm.DB) (Player, error) {

	var player Player
	err := db.
		Preload("User").
		Preload("PlayerStocks").
		Preload("Game").
		Preload("PlayerStocks.GameStock.Stock").
		Where("id = ?", playerID).
		First(&player).Error

	return player, err
}

func PlayerLeft(player *Player, db *gorm.DB) error {

	player.Active = false
	err := db.Save(player).Error

	if err != nil {
		fmt.Println("could not update player to inactive")
		return err
	}

	return err
}

// sort array of players by ID bubble sort
func SortPlayers(players []Player) []Player {

	fmt.Println("sorting players", len(players))

	for i := 0; i < len(players); i++ {
		fmt.Println("active?", players[i].User.Name, players[i].Active)

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

// sort array of playerstocks by gamestock.stock.variation bubble sort
func (player Player) SortedPlayerStocks() []PlayerStock {

	player_stocks := player.PlayerStocks
	for i := 0; i < len(player_stocks); i++ {
		for j := 0; j < len(player_stocks)-i-1; j++ {
			if player_stocks[j].GameStock.Stock.Variation > player_stocks[j+1].GameStock.Stock.Variation {
				player_stocks[j], player_stocks[j+1] = player_stocks[j+1], player_stocks[j]
			}
		}
	}

	return player_stocks
}
