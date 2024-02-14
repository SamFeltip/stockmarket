package models

import (
	"fmt"

	gorm "gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null;unique"`
	StartingValue float64
	ImagePath     string  `gorm:"not null;unique"`
	Variation     float64 // +/- maximum value of variation (0.50 increments)
}

type GameStock struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	GameID       string `gorm:"not null"`
	StockID      uint   `gorm:"not null"`
	Stock        Stock
	PlayerStocks []PlayerStock
	DirectorID   *uint
	Director     *Player `gorm:"foreignKey:DirectorID"`
	Game         Game
	Value        float64
}

type PlayerStock struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	PlayerID    uint
	Player      Player
	GameStockID uint
	GameStock   GameStock
	Quantity    int
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
		Where("id = ?", playerStockID).First(&playerStock).Error

	return playerStock, err
}

func CreateGameStocks(gameID string, db *gorm.DB) ([]GameStock, error) {
	fmt.Println("creating game stocks")

	var stocks []Stock
	db.Find(&stocks)

	fmt.Println("stocks.len: ", len(stocks))

	var game_stocks []GameStock

	for _, stock := range stocks {
		fmt.Println("stock iteration:", stock.ID, ":", gameID)
		game_stock := GameStock{
			GameID:   gameID,
			StockID:  stock.ID,
			Value:    stock.StartingValue,
			Director: nil,
		}

		err := db.Create(&game_stock).Error

		if err != nil {
			fmt.Println("failed to create game stock", err)
			return nil, err
		}

		game_stocks = append(game_stocks, game_stock)
	}

	return game_stocks, nil
}
