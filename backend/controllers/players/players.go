package controllers

import (
	"fmt"
	"stockmarket/models"
	templates "stockmarket/templates/players"

	"github.com/a-h/templ"
	"gorm.io/gorm"
)

func Show(playerID uint, currentPlayerID uint, db *gorm.DB) templ.Component {

	playerStockDisplays, err := models.GetPlayerStockDisplays(playerID, db)

	if err != nil {
		fmt.Println("error loading player:", err)
		pageComponent := templates.NoPlayer(fmt.Errorf("could not find player stock display"))
		return pageComponent
	}

	insights := []models.InsightResult{}
	if playerID != currentPlayerID {
		fmt.Println("not showing insights for other players")
		pageComponent := templates.PlayerPortfolio(playerStockDisplays, insights)
		return pageComponent
	}

	fmt.Println("getting insights for current player:", playerID)
	err = db.Table("player_insights as pi").
		Select("i.value, i.description, s.name as stock_name, s.image_path as stock_image_path").
		Joins("inner join insights as i on i.id = pi.insight_id").
		Joins("inner join player_stocks as ps on ps.id = pi.player_stock_id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Joins("inner join stocks as s on s.id = gs.stock_id").
		Where("ps.player_id = ?", playerID).
		Order("s.Variation").
		Scan(&insights).Error

	if err != nil {
		fmt.Println("error loading insights:", err)
	}

	pageComponent := templates.PlayerPortfolio(playerStockDisplays, insights)
	return pageComponent
}
