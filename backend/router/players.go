package router

import (
	"context"
	"fmt"
	"stockmarket/database"
	"stockmarket/models"
	templates "stockmarket/templates/players"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePlayerRoutes() {

	r.GET("/players/show/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		db := database.GetDb()

		iduint64, err := strconv.ParseUint(idParam, 10, 64)

		if err != nil {
			fmt.Println("error converting id to uint", err)
			pageComponent := templates.NoPlayer(fmt.Errorf("error converting id to uint"))
			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)
			return
		}

		playerID := uint(iduint64)

		playerStockDisplays, err := models.GetPlayerStockDisplays(playerID, db)

		if err != nil {
			fmt.Println("error loading player:", err)
			pageComponent := templates.NoPlayer(fmt.Errorf("no user found in request context"))
			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)
			return
		}

		pageComponent := templates.PlayerPortfolio(playerStockDisplays)
		ctx := context.Background()
		pageComponent.Render(ctx, c.Writer)
	})
}
