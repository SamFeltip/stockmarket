package games

import (
	"stockmarket/models"
	templates "stockmarket/templates/games"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Show(c *gin.Context, db *gorm.DB) templ.Component {
	id := c.Param("id")

	var game models.Game
	db.First(&game, id)

	pageComponent := templates.Waiting(game)
	return pageComponent

}

func Create(c *gin.Context, db *gorm.DB) models.Game {
	name := c.PostForm("name")
	difficulty := c.PostForm("difficulty")

	game := models.Game{
		Name:       name,
		Difficulty: difficulty,
	}

	db.Create(&game)

	return game // passed into templates
}

func New(c *gin.Context, db *gorm.DB) {

}
