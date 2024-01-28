package users

import (
	"stockmarket/database"
	"stockmarket/models"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) []models.User {

	db := database.GetDb()

	// get all users from gorm
	var users []models.User
	db.Find(&users)

	return users // passed into templates
}

func Show(c *gin.Context) models.User {

	db := database.GetDb()

	id := c.Param("id")

	var user models.User
	db.First(&user, id)

	return user // passed into templates

}
