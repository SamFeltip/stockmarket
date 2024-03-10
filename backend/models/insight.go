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

type GameInsight struct {
	Description    string
	InsightValue   float64
	GameStockID    uint
	Name           string
	ImagePath      string
	GameStockValue float64
}

func GetGameInsights(gameID string, db *gorm.DB) ([]GameInsight, error) {

	var gameInsights []GameInsight
	err := db.Table("game_stocks as gs").
		Select("i.description, i.value as insight_value, gs.ID as game_stock_id, s.name, s.image_path, gs.value as game_stock_value").
		Joins("inner join stocks as s on s.id = gs.stock_id").
		Joins("inner join player_stocks as ps on ps.game_stock_id = gs.id").
		Joins("left join player_insights as pi on pi.player_stock_id = ps.id").
		Joins("left join insights as i on i.id = pi.insight_id").
		Where("gs.game_id = ?", gameID).
		Order("s.variation").
		Scan(&gameInsights).Error

	return gameInsights, err
}
