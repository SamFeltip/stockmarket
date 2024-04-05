package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Insight struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	StockID     uint
	Stock       Stock
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

// todo: reconsile gameInsight and InsightDisplay
type InsightDisplay struct {
	Description    string
	Value          float64
	GameStockID    uint
	StockName      string
	StockImagePath string
	StockDisplay   bool
}

type GameInsight struct {
	InsightID          uint
	Description        string
	InsightValue       float64
	GameStockID        uint
	Name               string
	ImagePath          string
	SecondaryImagePath string
	GameStockValue     float64
	UserName           string
}

func LoadSpecialPlayerInsights(playerID uint, db *gorm.DB) ([]InsightDisplay, error) {

	var insights []InsightDisplay

	fmt.Println("fetching special player insights")

	err := db.Table("player_insights as pi").
		Select("i.description, i.value, s.image_path as stock_image_path, s.name as stock_name, s.display as stock_display").
		Joins("inner join insights as i on i.id = pi.insight_id").
		Joins("inner join player_stocks as ps on ps.id = pi.player_stock_id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Joins("inner join stocks as s on s.id = gs.stock_id").
		Where("s.display = false AND ps.player_id = ?", playerID).Scan(&insights).
		Error

	return insights, err
}

func LoadSpecialInsights(gameID string, db *gorm.DB) ([]GameInsight, error) {
	insights := []GameInsight{}

	err := db.Table("player_insights as pi").
		Select("i.Id as insight_id, i.description, i.value as insight_value, gs.id as game_stock_id, s.name, s.image_path, COALESCE(s.secondary_image_path, u.profile_root) as secondary_image_path, u.name as user_name").
		Joins("inner join insights as i on i.id = pi.insight_id").
		Joins("inner join stocks as s on s.id = i.stock_id").
		Joins("inner join player_stocks as ps on ps.id = pi.player_stock_id").
		Joins("inner join players as p on p.id = ps.player_id").
		Joins("inner join users as u on u.id = p.user_id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Where("gs.game_id = ? AND s.display = false", gameID).Scan(&insights).Error

	return insights, err

}

func GetGameInsights(gameID string, db *gorm.DB) ([]GameInsight, error) {

	var gameInsights []GameInsight
	err := db.Table("game_stocks as gs").
		Select("i.description, i.value as insight_value, gs.ID as game_stock_id, s.name, s.image_path, gs.value as game_stock_value").
		Joins("inner join stocks as s on s.id = gs.stock_id").
		Joins("inner join player_stocks as ps on ps.game_stock_id = gs.id").
		Joins("left join player_insights as pi on pi.player_stock_id = ps.id").
		Joins("left join insights as i on i.id = pi.insight_id").
		Where("gs.game_id = ? AND s.display = true", gameID).
		Order("s.variation").
		Scan(&gameInsights).Error

	return gameInsights, err
}
