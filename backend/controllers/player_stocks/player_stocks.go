package player_stocks

import (
	"fmt"
	gameController "stockmarket/controllers/games"
	"stockmarket/models"
	gameTemplates "stockmarket/templates/games"

	"github.com/a-h/templ"
	"gorm.io/gorm"
)

func Edit(playerStockID uint, gameID string, quantityAdd int, multiplier int, db *gorm.DB) (templ.Component, error) {

	quantityChange := quantityAdd * multiplier

	_, err := models.NewFeedItem(quantityChange, playerStockID, db)

	if err != nil {
		fmt.Println("could not create new feed item", err)
		return gameTemplates.Error(err), err
	}

	playerStock := models.PlayerStock{}
	err = db.Where("id = ?", playerStockID).First(&playerStock).Error

	if err != nil {
		fmt.Println("could not get player stock", err)
		return gameTemplates.Error(err), err
	}

	newStockQuantity := playerStock.Quantity + quantityChange

	err = db.Model(models.PlayerStock{}).Where("id = ?", playerStockID).Update("quantity", newStockQuantity).Error

	if err != nil {
		fmt.Println("could not update player stock quantity", err)
		return gameTemplates.Error(err), err
	}

	playerChangeResult := struct {
		PlayerID       uint
		PlayerCash     int
		GameStockValue float64
	}{}

	err = db.Model(&playerChangeResult).Table("player_stocks as ps").
		Select("p.id as player_id, p.cash as player_cash, gs.value as game_stock_value").
		Joins("inner join players as p on p.id = ps.player_id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Where("ps.id = ?", playerStockID).
		First(&playerChangeResult).Error

	if err != nil {
		fmt.Println("could not get playerChangeResult", err)
		return gameTemplates.Error(err), err
	}

	newPlayerCash := float64(playerChangeResult.PlayerCash) - float64(quantityAdd)*playerChangeResult.GameStockValue

	err = db.Model(models.Player{}).Where("id = ?", playerChangeResult.PlayerID).Update("cash", newPlayerCash).Error

	if err != nil {
		fmt.Println("could not update player cash", err)
		return gameTemplates.Error(err), err
	}

	template, err := gameController.CheckForMarketClose(gameID, db)

	if err != nil {
		fmt.Println("could not check for market close", err)
		return gameTemplates.Error(err), err
	}

	if template != nil {
		fmt.Println("market closed")
		return template, nil
	}

	fmt.Println("market not closed")

	userID, err := models.UpdateCurrentUser(gameID, db)
	if err != nil {
		fmt.Println("could not update current player", err)

		return gameTemplates.Error(err), err
	}

	fmt.Println("updated game current user", userID)

	// get game loading template
	loadingComponent := gameTemplates.Loading()

	gameController.BroadcastUpdatePlayBoard(gameID)

	return loadingComponent, nil
}
