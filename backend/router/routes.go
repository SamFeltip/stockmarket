package router

import (
	"context"
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
	baseComponent := page.Base(title, pageComponent)
	baseComponent.Render(context.Background(), c.Writer)
}
