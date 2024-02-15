package models

import "gorm.io/gorm"

type Insight struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	StockID uint
	Stock   Stock
	// PlayerInsights []PlayerInsight (not a necessary field)
	Description string
	Value       float64
}

type PlayerInsight struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"`
	PlayerStockID uint
	PlayerStock   PlayerStock
	InsightID     uint
	Insight       Insight
}
