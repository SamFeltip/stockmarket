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
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateGameRoutes() {

	r.GET("/games/show/:id",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {

			db := database.GetDb()

			gameID := strings.ToLower(c.Param("id"))

			cu, exists := c.Get("user")

			if !exists {
				fmt.Println("no user found")
				pageComponent := templates.Error(fmt.Errorf("no user found in request context"))
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
			}

			current_user := cu.(models.User)

			pageComponent := controllers.Show(gameID, current_user, db)
			gameWrapper := templates.Base(pageComponent, gameID)

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

			code := strings.ToLower(c.PostForm("code"))
			periodCountStr := c.PostForm("difficulty")

			if code == "" || periodCountStr == "" {
				fmt.Println("no code or periodCount in form")
				pageComponent := templates.Create("no code or periodCount in form")
				RenderWithTemplate(pageComponent, "Create new game", c)
				return
			}

			periodCount, err := strconv.Atoi(periodCountStr)
			if err != nil {
				// handle error, e.g. return an error response
				fmt.Println("couldnt convert to int")
				pageComponent := templates.Create("couldnt convert to int")
				RenderWithTemplate(pageComponent, "Create new game", c)
				return
			}

			cu, _ := c.Get("user")
			current_user := cu.(models.User)

			if err != nil {
				fmt.Println("error creating game stocks, creating empty set:", err)
				pageComponent := templates.Create("error creating game stocks, creating empty set")
				RenderWithTemplate(pageComponent, "Create new game", c)
				return

			}

			db := database.GetDb()
			game, err := controllers.Create(code, periodCount, current_user, db)

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

			gameID := c.PostForm("gameID")

			cu, exists := c.Get("user")

			if !exists {
				fmt.Println("no user found")
				pageComponent := templates.Error(fmt.Errorf("no user found in request context"))
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
			}

			current_user := cu.(models.User)

			pageComponent, err := controllers.PlayAction(gameID, current_user, db)

			if err != nil {
				fmt.Println("error editing player stock", err)
				return
			}

			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)
		})

	r.POST("/api/games/next",
		func(c *gin.Context) { middleware.AuthIsPlaying(c) },
		func(c *gin.Context) {
			db := database.GetDb()
			gameID := c.PostForm("gameID")

			if gameID == "" {
				fmt.Println("no gameID in post request")
				pageComponent := templates.Error(fmt.Errorf("no gameID in post request"))
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			fmt.Println("form data gathered", "gameID:", gameID)

			pageComponent, err := controllers.NextPeriod(gameID, db)

			if err != nil {
				fmt.Println("error editing player stock", err)
				return
			}

			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)
		})

	r.GET("/games",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {

			// get all games
			pageComponent := controllers.Index(c)

			RenderWithTemplate(pageComponent, "Show games", c)

		})
}
