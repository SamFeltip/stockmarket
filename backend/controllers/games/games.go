package games

import (
	"fmt"
	"stockmarket/database"
	"stockmarket/models"
	templates "stockmarket/templates/games"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Show(c *gin.Context) templ.Component {
	db := database.GetDb()

	gameID := c.Param("id")

	game, err := models.GetGame(gameID, db)

	if err != nil {
		fmt.Println("error fetching game:", err)
		return templates.NoGame()
	}

	cu, _ := c.Get("user")
	current_user := cu.(models.User)

	player, err := current_user.SetActiveGame(game, db)

	if err != nil {
		fmt.Println("error setting active game:", err)
	}

	game.Players = append(game.Players, player)

	err = game.BroadcastNewPlayer(player)

	if err != nil {
		fmt.Println("error broadcasting new player:", err)
		return templates.NoGame()
	}

	return templates.IngamePage(game)

}

func Create(c *gin.Context) models.Game {
	db := database.GetDb()

	code := c.PostForm("code")
	difficultyStr := c.PostForm("difficulty")

	difficulty, err := strconv.Atoi(difficultyStr)
	if err != nil {
		// handle error, e.g. return an error response
		fmt.Println("couldnt convert to int")
	}
	cu, _ := c.Get("user")
	current_user := cu.(models.User)

	game := models.Game{
		ID:          code,
		Difficulty:  difficulty,
		Status:      "waiting",
		CurrentUser: current_user,
	}

	db.Create(&game)

	return game // passed into templates
}

func New(c *gin.Context, db *gorm.DB) {

}

func Index(c *gin.Context) templ.Component {
	db := database.GetDb()

	// get all games from gorm
	var games []models.Game
	db.Find(&games)

	pageComponent := templates.Index(games)

	return pageComponent // passed into templates
}
