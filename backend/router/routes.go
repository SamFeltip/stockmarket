package router

import (
	"context"
	"net/http"
	"stockmarket/models"
	page "stockmarket/templates"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLFiles("sockets.html")

	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/sockets", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sockets.html", nil)
	})

	CreateAuthRoutes(db, r)

	CreateWebsocketRoutes(db, r)
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
