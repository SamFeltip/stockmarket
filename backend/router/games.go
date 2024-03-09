package router

import (
	"context"
	"fmt"
	"net/http"
	controllers "stockmarket/controllers/games"
	"stockmarket/database"
	"stockmarket/middleware"
	models "stockmarket/models"
	templates "stockmarket/templates/games"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGameRoutes() {

	r.GET("/games/show/:id",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {

			db := database.GetDb()

			gameID := c.Param("id")

			game, err := models.LoadGame(gameID, db)

			if err != nil {
				fmt.Println("error fetching game:", err)
				gameWrapper := templates.NoGame()
				RenderWithTemplate(gameWrapper, "Game - id", c)
				return
			}

			c.Set("game", game)

			pageComponent := controllers.Show(db, c)
			gameWrapper := templates.Base(pageComponent, game)

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

			// log form in context (form contains gameID and periodCount)
			c.Request.ParseForm()
			fmt.Println(c.Request.Form["gameID"])
			fmt.Println(c.Request.Form["game-length"])

			gameID := c.Request.Form["gameID"][0]
			periodCountStr := c.Request.Form["game-length"][0]

			periodCount, err := strconv.Atoi(periodCountStr)
			if err != nil {
				fmt.Println("could not convert periodCount to int", periodCountStr)
				pageComponent := templates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			baseComponent, err := controllers.UpdateGamePeriodCount(gameID, periodCount)
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

			// log form in context (form contains gameID and periodCount)
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

	r.POST("/api/games/action",
		func(c *gin.Context) { middleware.AuthCurrentPlayer(c) },
		func(c *gin.Context) {
			db := database.GetDb()
			gameAction := c.PostForm("game_action")

			if gameAction == "" {
				fmt.Println("no gameID or play action in post request")
				pageComponent := templates.Error(fmt.Errorf("no gameID or play action in post request"))
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			fmt.Println("form data gathered", "gameAction:", gameAction)

			pageComponent, err := controllers.PlayAction(c, db)
			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)

			if err != nil {
				fmt.Println("error editing player stock, don't broadcast", err)
				return
			}

			cg, _ := c.Get("game")
			game := cg.(models.Game)

			controllers.BroadcastUpdateBoard(game)
		})

	r.GET("/games",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {

			// get all games
			pageComponent := controllers.Index(c)

			RenderWithTemplate(pageComponent, "Show games", c)

		})
}
