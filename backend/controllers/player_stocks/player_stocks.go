package player_stocks

import (
	"fmt"
	"stockmarket/models"
	gameTemplates "stockmarket/templates/games"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Edit(c *gin.Context, db *gorm.DB) (templ.Component, error) {
	playerStockID := c.PostForm("PlayerStockID")
	playerStockQuantityAdd := c.PostForm("PlayerStockQuantityAdd")

	playerStock, err := models.GetPlayerStock(playerStockID, db)

	if err != nil {
		fmt.Println("could not find player stock", err)
		return gameTemplates.Error(err), err
	}

	// parse QuantityAdd to int and add to player stock . quantity
	quantityAdd, err := strconv.Atoi(playerStockQuantityAdd)
	if err != nil {
		fmt.Println("could not parse new quantity to int", err)
		return gameTemplates.Error(err), err
	}

	playerStock.Quantity += quantityAdd

	db.Save(&playerStock)

	cg, exists := c.Get("game")

	if !exists {
		fmt.Println("could not get game from context", err)
		return gameTemplates.Error(err), err
	}

	game := cg.(models.Game)

	_, err = models.NewFeedItem(playerStock, game, quantityAdd, db)

	if err != nil {
		fmt.Println("could not create new feed item", err)
		return gameTemplates.Error(err), err
	}

	err = game.UpdateCurrentUser(db)

	if err != nil {
		fmt.Println("could not update current player", err)

		return gameTemplates.Error(err), err
	}

	fmt.Println("updated game current user", game.CurrentUser.Name)

	c.Set("game", game)

	// get game loading template
	loadingComponent := gameTemplates.Loading()

	return loadingComponent, nil
}
