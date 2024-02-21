package router

import (
	"context"
	"stockmarket/database"
	"stockmarket/middleware"
	models "stockmarket/models"
	templates "stockmarket/templates/player_stocks"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func CreatePlayerStockRoutes() {

	r.GET("/player_stocks/show/:playerStockID",
		func(c *gin.Context) { middleware.RequireAuth(c) },
		func(c *gin.Context) {
			db := database.GetDb()

			playerStockID := c.Param("playerStockID")
			playerStock, err := models.GetPlayerStock(playerStockID, db)

			var pageComponent templ.Component
			if err != nil {
				pageComponent = templates.NoPlayerStock()
			} else {
				pageComponent = templates.Show(playerStock)
			}

			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)

		})

	r.GET("/player_stocks/preview/:playerStockID",
		func(c *gin.Context) { middleware.RequireAuth(c) },
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

}
