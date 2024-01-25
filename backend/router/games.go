package router

import (
	"fmt"
	"net/http"
	controllers "stockmarket/controllers/games"
	"stockmarket/middleware"
	"stockmarket/models"
	templates "stockmarket/templates/games"
	"strconv"

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
			RenderWithTemplate(pageComponent, "Create new game", c)

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

	r.POST(
		"/games/cmd/difficulty",
		func(c *gin.Context) { middleware.AuthCurrentPlayer(c, db) },
		func(c *gin.Context) {
			// log form in context (form contains gameID and difficulty)
			c.Request.ParseForm()
			fmt.Println(c.Request.Form["gameID"])
			fmt.Println(c.Request.Form["game-length"])

			gameID := c.Request.Form["gameID"][0]
			difficultyStr := c.Request.Form["game-length"][0]

			difficulty, err := strconv.Atoi(difficultyStr)
			if err != nil {
				fmt.Println("could not convert difficulty to int", difficultyStr)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "could not convert difficulty to int",
				})
				return
			}

			// Assuming you have a *gorm.DB object named db
			err = db.Model(models.Game{}).Where("id = lower(?)", gameID).Update("difficulty", difficulty).Error

			if err != nil {
				fmt.Println("could not update game difficulty")
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "could not update game difficulty",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"error":   "",
				"message": "difficulty updated",
			})

		},
	)

	r.GET(
		"/games",
		func(c *gin.Context) { middleware.RequireAuth(c, db) },
		func(c *gin.Context) {

			// get all games
			pageComponent := controllers.Index(c, db)

			RenderWithTemplate(pageComponent, "Show games", c)

		})
}
