package games

import (
	"fmt"
	"stockmarket/models"
	templates "stockmarket/templates/games"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Show(c *gin.Context, db *gorm.DB) templ.Component {
	id := c.Param("id")

	var game models.Game
	err := db.Where("lower(id) = lower(?)", id).First(&game).Error

	if err != nil {
		fmt.Println("error fetching game:", err)
		pageComponent := templates.NoGame()
		return pageComponent
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
