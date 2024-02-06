package models

import gorm "gorm.io/gorm"

type Stock struct {
	gorm.Model
	ID            uint
	Name          string
	StartingValue float64
	ImagePath     string
	Variation     float64 // +/- maximum value of variation (0.50 increments)
}

type GameStock struct {
	gorm.Model
	ID      uint
	GameID  string
	StockID string
	Stock   Stock
	Game    Game
	Value   float64
}

type PlayerStock struct {
	gorm.Model
	ID          uint
	PlayerID    string
	GameStockID uint
	GameStock   GameStock
	Player      Player
	Quantity    int
}
