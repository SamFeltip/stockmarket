package router

import (
	"fmt"
	"net/http"
	controllers "stockmarket/controllers/games"
	templates "stockmarket/templates/games"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateGameRoutes(db *gorm.DB, r *gin.Engine) {

	r.GET("/games/show/:id", func(c *gin.Context) {

		pageComponent := controllers.Show(c, db)

		RenderWithTemplate(pageComponent, "Game - id", c)

	})

	r.GET("/games/new", func(c *gin.Context) {

		pageComponent := templates.New()
		RenderWithTemplate(pageComponent, "Signup", c)

	})

	r.POST("/games/new", func(c *gin.Context) {

		game := controllers.Create(c, db)

		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/games/show/%d?status=waiting", game.ID))

	})
}
