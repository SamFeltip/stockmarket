package router

import (
	"context"
	"fmt"
	controllers "stockmarket/controllers/players"
	"stockmarket/database"
	"stockmarket/middleware"
	"stockmarket/models"
	templates "stockmarket/templates/players"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePlayerRoutes() {

	r.GET("/players/show/:id",
		func(c *gin.Context) { middleware.AuthIsPlaying(c) },
		func(c *gin.Context) {
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

			cp, exists := c.Get("player")

			if !exists {
				fmt.Println("no player found in request context")
				pageComponent := templates.NoPlayer(fmt.Errorf("no player found in request context"))
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			currentPlayer := cp.(models.Player)

			pageComponent := controllers.Show(playerID, currentPlayer.ID, db)

			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)
		})
}
