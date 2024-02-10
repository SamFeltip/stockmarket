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

}
