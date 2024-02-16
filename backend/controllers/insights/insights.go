package insights

import (
	"stockmarket/database"
	"stockmarket/models"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) []models.Insight {
	db := database.GetDb()

	// get all games from gorm
	var insights []models.Insight
	db.Preload("Stock").Find(&insights)

	return insights // passed into templates
}
