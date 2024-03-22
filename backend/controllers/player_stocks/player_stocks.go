package player_stocks

import (
	"fmt"
	gameController "stockmarket/controllers/games"
	"stockmarket/models"
	gameTemplates "stockmarket/templates/games"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Edit(c *gin.Context, db *gorm.DB) (templ.Component, error) {
	playerStockIDString := c.PostForm("PlayerStockID")
	playerStockQuantityAdd := c.PostForm("PlayerStockQuantityAdd")
	gameID := c.PostForm("gameID")

	playerStockID64, err := strconv.ParseUint(playerStockIDString, 10, 32)

	if err != nil {
		fmt.Println("could not parse id", err)
		return gameTemplates.Error(err), err
	}

	playerStockID := uint(playerStockID64)

	// parse QuantityAdd to int and add to player stock . quantity
	quantityAdd, err := strconv.Atoi(playerStockQuantityAdd)
	if err != nil {
		fmt.Println("could not parse new quantity to int", err)
		return gameTemplates.Error(err), err
	}

	mode := c.PostForm("mode")

	multiplier, err := strconv.Atoi(mode)
	if err != nil {
		fmt.Println("could not parse mode to int", err)
		return gameTemplates.Error(err), err
	}

	// completed form validation

	quantityChange := quantityAdd * multiplier

	_, err = models.NewFeedItem(quantityChange, playerStockID, db)

	if err != nil {
		fmt.Println("could not create new feed item", err)
		return gameTemplates.Error(err), err
	}

	// update playerstock with playerstockid to quantiy of given
	// update player cash with quantity of given

	// Assuming playerStockID and quantityChange are defined
	err = db.Model(models.PlayerStock{}).Where("id = ?", playerStockID).Update("quantity", quantityChange).Error

	if err != nil {
		fmt.Println("could not update player stock quantity", err)
		return gameTemplates.Error(err), err
	}

	playerChangeResult := struct {
		PlayerID       uint
		PlayerCash     int
		GameStockValue int
	}{}

	err = db.Model(&playerChangeResult).Table("player_stocks as ps").
		Select("p.id as player_id, p.cash as player_cash, gs.value as game_stock_value").
		Joins("inner join players as p on p.id = ps.player_id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").Error

	newPlayerCash := playerChangeResult.PlayerCash - quantityAdd*int(playerChangeResult.GameStockValue)

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
