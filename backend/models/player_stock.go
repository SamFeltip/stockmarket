package models

import (
	"fmt"

	"gorm.io/gorm"
)

type PlayerStockPreview struct {
	TotalInsight float64
	StockValue   float64
	StockName    string
	StockImg     string
}

type PlayerStock struct {
	gorm.Model
	ID             uint `gorm:"primaryKey"`
	PlayerID       uint
	Player         Player
	GameStockID    uint
	GameStock      GameStock
	PlayerInsights []PlayerInsight `gorm:"constraint:OnDelete:CASCADE"`
	Quantity       int
}

func GetPlayerStock(playerStockID string, db *gorm.DB) (PlayerStock, error) {
	var playerStock PlayerStock

	// player_stock.GameStock.Stock.Name

	err := db.
		Preload("GameStock").
		Preload("GameStock.Stock").
		Preload("GameStock.PlayerStocks").
		Preload("GameStock.PlayerStocks.Player").
		Preload("GameStock.PlayerStocks.Player.User").
		Preload("Player").
		Preload("Player.User").
		Preload("PlayerInsights").
		Preload("PlayerInsights.Insight").
		Where("id = ?", playerStockID).First(&playerStock).Error

	return playerStock, err
}

func (player_stock PlayerStock) TotalInsight() float64 {

	fmt.Println("player_stock stock: ", player_stock.GameStock.Stock.Name)

	var total float64 = 0

	for _, player_insight := range player_stock.PlayerInsights {
		total += player_insight.Insight.Value
	}

	return total
}

func (player_stock PlayerStock) Value() float64 {
	return float64(player_stock.Quantity) * player_stock.GameStock.Value
}
