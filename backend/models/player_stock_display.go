package models

import (
	"fmt"

	"gorm.io/gorm"
)

type PlayerStockDisplay struct {
	ID                  uint
	GameID              string
	TotalInsight        float64
	GameStockValue      float64
	PlayerStockQuantity int
	StockName           string
	StockImagePath      string
}

func GetPlayerStockDisplays(playerID uint, db *gorm.DB) ([]PlayerStockDisplay, error) {
	var playerStocksResult = []PlayerStockDisplay{}

	err := db.Table("player_stocks as ps").
		Select("ps.ID, gs.game_id, gs.value as game_stock_value, ps.quantity as player_stock_quantity, s.name as stock_name, s.image_path as stock_image_path, COALESCE(sum(i.value), 0) as total_insight").
		Joins("left join player_insights as pi on pi.player_stock_id = ps.id").
		Joins("left join insights as i on i.id = pi.insight_id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Joins("inner join stocks as s on s.id = gs.stock_id").
		Where("ps.player_id = ?", playerID).
		Group("ps.ID, gs.game_id, gs.value, s.name, s.image_path, s.variation").
		Order("s.variation").
		Scan(&playerStocksResult).Error

	if err != nil {
		fmt.Println("could not load player stocks", err)
		return nil, err
	}

	if len(playerStocksResult) == 0 {
		fmt.Println("no player stocks found for this player")
		return nil, gorm.ErrRecordNotFound
	}

	return playerStocksResult, nil
}
