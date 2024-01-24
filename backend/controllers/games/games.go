package games

import (
	"fmt"
	"stockmarket/models"
	templates "stockmarket/templates/games"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Show(c *gin.Context, db *gorm.DB) templ.Component {
	gameID := c.Param("id")

	var game models.Game
	err := db.Model(&game).Preload(clause.Associations).Where("lower(id) = lower(?)", gameID).First(&game).Error

	if err != nil {
		fmt.Println("error fetching game:", err)
		pageComponent := templates.NoGame()
		return pageComponent
	}

	cu, _ := c.Get("user")
	current_user := cu.(models.User)

	player, err := models.GetPlayer(game, current_user, db)

	if err != nil {
		fmt.Println("error fetching player:", err)

		if err == gorm.ErrRecordNotFound {
			player = models.Player{
				GameID: game.ID,
				UserID: current_user.ID,
			}
			db.Create(&player)
			game.Players = append(game.Players, player)
		}
	}

	if game.Status == "playing" {
		pageComponent := templates.Playing(game)
		return pageComponent
	}

	if game.Status == "showing" {
		pageComponent := templates.Showing(game)
		return pageComponent
	}

	pageComponent := templates.Waiting(game)
	return pageComponent

}

func Create(c *gin.Context, db *gorm.DB) models.Game {
	code := c.PostForm("code")
	difficultyStr := c.PostForm("difficulty")

	difficulty, err := strconv.Atoi(difficultyStr)
	if err != nil {
		// handle error, e.g. return an error response
		fmt.Println("couldnt convert to int")
	}

	game := models.Game{
		ID:         code,
		Difficulty: difficulty,
		Status:     "waiting",
	}

	db.Create(&game)

	return game // passed into templates
}

func New(c *gin.Context, db *gorm.DB) {

}

func Index(c *gin.Context, db *gorm.DB) templ.Component {

	// get all games from gorm
	var games []models.Game
	db.Find(&games)

	pageComponent := templates.Index(games)

	return pageComponent // passed into templates
}