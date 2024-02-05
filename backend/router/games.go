package router

import (
	"context"
	"fmt"
	"net/http"
	controllers "stockmarket/controllers/games"
	"stockmarket/database"
	"stockmarket/middleware"
	"stockmarket/models"
	templates "stockmarket/templates/games"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGameRoutes() {

	r.GET(
		"/games/show/:id",
		func(c *gin.Context) { middleware.RequireAuth(c) },
		func(c *gin.Context) {

			pageComponent := controllers.Show(c)

			RenderWithTemplate(pageComponent, "Game - id", c)

		})

	r.GET(
		"/games/new",
		func(c *gin.Context) { middleware.RequireAuth(c) },
		func(c *gin.Context) {

			pageComponent := templates.Create()
			RenderWithTemplate(pageComponent, "Create new game", c)

		})

	r.POST(
		"/games/new",
		func(c *gin.Context) { middleware.RequireAuth(c) },
		func(c *gin.Context) {

			game := controllers.Create(c)

			show_url := fmt.Sprintf("/games/show/%s", game.ID)

			fmt.Println(show_url)

			c.Redirect(http.StatusMovedPermanently, show_url)

		})

	r.POST(
		"api/games/difficulty",
		func(c *gin.Context) { middleware.AuthCurrentPlayer(c) },
		func(c *gin.Context) {

			db := database.GetDb()
			errMsg := ""

			// log form in context (form contains gameID and difficulty)
			c.Request.ParseForm()
			fmt.Println(c.Request.Form["gameID"])
			fmt.Println(c.Request.Form["game-length"])

			gameID := c.Request.Form["gameID"][0]
			difficultyStr := c.Request.Form["game-length"][0]

			difficulty, err := strconv.Atoi(difficultyStr)
			if err != nil {
				fmt.Println("could not convert difficulty to int", difficultyStr)
				errMsg = "could not convert difficulty to int"
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": errMsg,
				})
			}

			game := models.Game{}
			err = db.Model(&game).Where("id = lower(?)", gameID).Update("difficulty", difficulty).Error

			if err != nil {
				fmt.Println("could not update game difficulty")
				errMsg = "could not update game difficulty"
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": errMsg,
				})
			}

			err = db.Model(&game).Where("lower(id) = lower(?)", gameID).First(&game).Error

			if err != nil {
				fmt.Printf("Error reloading game: %v", err)
				errMsg = "Error reloading game"
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": errMsg,
				})
			}

			err = controllers.BroadcastUpdateDifficulty(game)

			if err != nil {
				fmt.Println("could not broadcast difficulty update")
				errMsg = "could not broadcast difficulty update"
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": errMsg,
				})
			}

			// return DifficultyOptions(game) templ
			ctx := context.Background()

			baseComponent := templates.DifficultyOptions(game, errMsg)
			baseComponent.Render(ctx, c.Writer)

		},
	)

	r.GET(
		"/games",
		func(c *gin.Context) { middleware.RequireAuth(c) },
		func(c *gin.Context) {

			// get all games
			pageComponent := controllers.Index(c)

			RenderWithTemplate(pageComponent, "Show games", c)

		})
}
