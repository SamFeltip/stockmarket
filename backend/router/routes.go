package router

import (
	"context"
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

	return r
}

func RenderWithTemplate(pageComponent templ.Component, title string, c *gin.Context) {

	cu, _ := c.Get("user")

	if cu == nil {
		cu = models.User{}
	}

	user := cu.(models.User)

	ctx := context.WithValue(context.Background(), page.CurrentUser, user)

	baseComponent := page.Base(title, pageComponent, c)
	baseComponent.Render(ctx, c.Writer)
}
