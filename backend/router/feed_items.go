package router

import (
	"context"
	"fmt"
	"stockmarket/database"
	"stockmarket/middleware"
	"stockmarket/models"
	templates "stockmarket/templates/feed_items"

	"github.com/gin-gonic/gin"
)

func CreateFeedItemRoutes() {

	r.GET("/feed_items/show/:game_id",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {

			db := database.GetDb()

			gameID := c.Param("game_id")

			feed, err := models.LoadFeedItems(gameID, db)

			if err != nil {
				fmt.Println("error fetching feed items:", err)
				feedTemplate := templates.NoFeed()
				ctx := context.Background()
				feedTemplate.Render(ctx, c.Writer)
				return
			}

			feedTemplate := templates.Feed(feed)
			ctx := context.Background()
			feedTemplate.Render(ctx, c.Writer)
			return

		})
}
