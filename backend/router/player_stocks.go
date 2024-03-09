package router

import (
	"context"
	"fmt"
	gamecontrollers "stockmarket/controllers/games"
	controllers "stockmarket/controllers/player_stocks"
	"stockmarket/database"
	"stockmarket/middleware"
	models "stockmarket/models"
	gameTemplates "stockmarket/templates/games"
	templates "stockmarket/templates/player_stocks"

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
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return

			}

			currentPlayer := playerStock.Player

			game := models.Game{}
			err = db.Where("id = ?", playerStock.GameStock.GameID).First(&game).Error

			if err != nil {
				fmt.Println("error fetching game for player stock", err)
				pageComponent = gameTemplates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			isCurrentPlayer := currentPlayer.User.ID == game.CurrentUserID
			pageComponent = templates.Show(playerStock, isCurrentPlayer)

			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)

		})

	r.GET("/player_stocks/preview/:playerStockID",
		func(c *gin.Context) { middleware.AuthIsPlaying(c) },
		func(c *gin.Context) {
			db := database.GetDb()

			// get a player stock for the game stock and player
			player_stock_id := c.Param("playerStockID")

			var playerStockPreview models.PlayerStockPreview

			db.Table("player_stocks as ps").
				Select("sum(i.value) as total_insight, gs.value as stock_value, s.name as stock_name, s.image_path as stock_img").
				Joins("inner join player_insights as pi on pi.player_stock_id = ps.id").
				Joins("inner join insights as i on i.id = pi.insight_id").
				Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
				Joins("inner join stocks as s on s.id = gs.stock_id").
				Where("ps.id = ?", player_stock_id).
				Group("stock_value, stock_name, stock_img").
				Scan(&playerStockPreview)

			pageComponent := templates.PlayerStockPreview(playerStockPreview)

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

			fmt.Println("form data gathered: playerStockID:", playerStockID, "playerStockQuantityAdd:", playerStockQuantityAdd)

			pageComponent, err := controllers.Edit(c, db)
			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)

			if err != nil {
				fmt.Println("error editing player stock, don't broadcast", err)
				return
			}

			cg, _ := c.Get("game")
			game := cg.(models.Game)

			gamecontrollers.BroadcastUpdatePlayBoard(game)

		})
}
