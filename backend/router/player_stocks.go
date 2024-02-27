package router

import (
	"context"
	"fmt"
	gamecontrollers "stockmarket/controllers/games"
	"stockmarket/database"
	"stockmarket/middleware"
	models "stockmarket/models"
	gameTemplates "stockmarket/templates/games"
	templates "stockmarket/templates/player_stocks"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func CreatePlayerStockRoutes() {

	r.GET("/player_stocks/show/:playerStockID",
		func(c *gin.Context) { middleware.AuthIsPlaying(c) },
		func(c *gin.Context) {
			db := database.GetDb()

			playerStockID := c.Param("playerStockID")
			playerStock, err := models.GetPlayerStock(playerStockID, db)

			var pageComponent templ.Component
			if err != nil {
				pageComponent = templates.NoPlayerStock()
			} else {

				cp, _ := c.Get("player")
				cg, _ := c.Get("game")

				currentPlayer := cp.(models.Player)
				game := cg.(models.Game)

				pageComponent = templates.Show(playerStock, currentPlayer.User.ID == game.CurrentUserID)
			}

			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)

		})

	r.GET("/player_stocks/preview/:playerStockID",
		func(c *gin.Context) { middleware.AuthIsPlaying(c) },
		func(c *gin.Context) {
			db := database.GetDb()

			// get a player stock for the game stock and player
			player_stock_id := c.Param("playerStockID")

			player_stock := models.PlayerStock{}
			db.
				Preload("GameStock.Stock").
				Preload("PlayerInsights.Insight").
				Where("id = ?", player_stock_id).
				First(&player_stock)

			pageComponent := templates.PlayerStockPreview(player_stock)

			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)

		})

	r.POST("/player_stocks/edit",
		func(c *gin.Context) { middleware.AuthCurrentPlayer(c) },
		func(c *gin.Context) {
			db := database.GetDb()
			playerStockID := c.PostForm("PlayerStockID")
			playerStockQuantityAdd := c.PostForm("PlayerStockQuantityAdd")

			if playerStockID == "" || playerStockQuantityAdd == "" {
				fmt.Println("no playerStockID or playerStockQuantityAdd in form")
				pageComponent := gameTemplates.Error(fmt.Errorf("no playerStockID or playerStockQuantityAdd"))
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			fmt.Println("playerStockID", playerStockID)

			playerStock, err := models.GetPlayerStock(playerStockID, db)

			if err != nil {
				fmt.Println("could not find player stock", err)

				pageComponent := gameTemplates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			// parse QuantityAdd to int and add to player stock . quantity
			quantityAdd, err := strconv.Atoi(playerStockQuantityAdd)
			if err != nil {
				fmt.Println("could not parse new quantity to int", err)

				pageComponent := gameTemplates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			playerStock.Quantity += quantityAdd

			db.Save(&playerStock)

			cg, exists := c.Get("game")

			if !exists {
				fmt.Println("could not get game from context", err)
				pageComponent := gameTemplates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			game := cg.(models.Game)

			err = game.UpdateCurrentUser(db)

			if err != nil {
				fmt.Println("could not update current player", err)
				pageComponent := gameTemplates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			// get game loading template
			loadingComponent := gameTemplates.Loading()
			ctx := context.Background()
			loadingComponent.Render(ctx, c.Writer)

			gamecontrollers.BroadcastUpdateBoard(game)

		})
}
