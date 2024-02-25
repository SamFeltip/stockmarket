package router

import (
	"context"
	"fmt"
	"net/http"
	controllers "stockmarket/controllers/games"
	"stockmarket/middleware"
	templates "stockmarket/templates/games"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGameRoutes() {

	r.GET("/games/show/:id",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {

			pageComponent := controllers.Show(c)
			gameWrapper := templates.Base(pageComponent)
			RenderWithTemplate(gameWrapper, "Game - id", c)

		})

	r.GET("/games/new",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {

			pageComponent := templates.Create("")
			RenderWithTemplate(pageComponent, "Create new game", c)

		})

	r.POST("/games/new",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {

			game, err := controllers.Create(c)

			c.Set("game", game)

			if err != nil {
				fmt.Println("error creating game:", err)
				pageComponent := templates.Create(err.Error())
				RenderWithTemplate(pageComponent, "Create new game", c)
				return
			}

			show_url := fmt.Sprintf("/games/show/%s", game.ID)

			fmt.Println(show_url)

			c.Redirect(http.StatusMovedPermanently, show_url)

		})

	r.POST("/api/games/difficulty",
		func(c *gin.Context) { middleware.AuthCurrentPlayer(c) },
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
				pageComponent := templates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			baseComponent, err := controllers.UpdateGameDifficulty(gameID, difficulty)
			if err != nil {
				pageComponent := templates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)

				return
			}

			ctx := context.Background()
			baseComponent.Render(ctx, c.Writer)

		},
	)

	r.POST("/api/games/start",
		func(c *gin.Context) { middleware.AuthCurrentPlayer(c) },
		func(c *gin.Context) {
			fmt.Println("start game")

			// log form in context (form contains gameID and difficulty)
			c.Request.ParseForm()
			gameID := c.Request.Form["gameID"][0]

			baseComponent, err := controllers.StartGame(gameID)

			if err != nil {
				pageComponent := templates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)

				return
			}

			ctx := context.Background()
			baseComponent.Render(ctx, c.Writer)
		})

	r.GET("/games",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {

			// get all games
			pageComponent := controllers.Index(c)

			RenderWithTemplate(pageComponent, "Show games", c)

		})
}
