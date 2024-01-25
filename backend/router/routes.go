package router

import (
	"context"
	"stockmarket/models"
	page "stockmarket/templates"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	CreateAuthRoutes(db, r)

	CreatePageRoutes(db, r)
	CreateUserRoutes(db, r)
	CreateGameRoutes(db, r)

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
