package models

import (
	"gorm.io/gorm"
)

type PlayerStock struct {
	gorm.Model
	ID             uint `gorm:"primaryKey"`
	PlayerID       uint
	Player         Player
	GameStockID    uint
	GameStock      GameStock
	PlayerInsights []PlayerInsight
	Quantity       int
}

func (player_stock PlayerStock) TotalInsight() float64 {
	var total float64 = 0

	for _, player_insight := range player_stock.PlayerInsights {
		total += player_insight.Insight.Value
	}

	return total
}

func (player_stock PlayerStock) Value() float64 {
	return float64(player_stock.Quantity) * player_stock.GameStock.Value
}
