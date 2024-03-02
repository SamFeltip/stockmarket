package router

import (
	"context"
	"fmt"
	"stockmarket/models"
	page "stockmarket/templates"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func SetupRoutes() *gin.Engine {
	r = gin.Default()

	CreateAuthRoutes()

	CreateWebsocketRoutes()
	CreatePageRoutes()
	CreateUserRoutes()
	CreateGameRoutes()
	CreatePlayerStockRoutes()
	CreateInsightRoutes()
	CreateFeedItemRoutes()

	return r
}

func RenderWithTemplate(pageComponent templ.Component, title string, c *gin.Context) {

	cu, _ := c.Get("user")

	if cu == nil {
		cu = models.User{}
	}

	user := cu.(models.User)

	cg, _ := c.Get("game")

	if cg == nil {
		fmt.Println("no game found in context")
		cg = models.Game{}
	}

	game := cg.(models.Game)
	fmt.Println("^^^", game.CurrentUser.ID)

	ctx := context.Background()
	ctx = context.WithValue(ctx, page.CurrentUser, user)
	ctx = context.WithValue(ctx, page.CurrentGame, game)

	fmt.Println("rendering template with base:", title)

	baseComponent := page.Base(title, pageComponent, c)
	baseComponent.Render(ctx, c.Writer)
}
