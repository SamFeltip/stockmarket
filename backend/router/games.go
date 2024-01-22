package router

import (
	"fmt"
	"net/http"
	controllers "stockmarket/controllers/games"
	"stockmarket/middleware"
	templates "stockmarket/templates/games"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateGameRoutes(db *gorm.DB, r *gin.Engine) {

	r.GET(
		"/games/show/:id",
		func(c *gin.Context) { middleware.RequireAuth(c, db) },
		func(c *gin.Context) {

			pageComponent := controllers.Show(c, db)

			RenderWithTemplate(pageComponent, "Game - id", c)

		})

	r.GET(
		"/games/new",
		func(c *gin.Context) { middleware.RequireAuth(c, db) },
		func(c *gin.Context) {

			pageComponent := templates.Create()
			RenderWithTemplate(pageComponent, "Signup", c)

		})

	r.POST(
		"/games/new",
		func(c *gin.Context) { middleware.RequireAuth(c, db) },
		func(c *gin.Context) {

			game := controllers.Create(c, db)

			show_url := fmt.Sprintf("/games/show/%s", game.ID)

			fmt.Println(show_url)

			c.Redirect(http.StatusMovedPermanently, show_url)

		})
}
