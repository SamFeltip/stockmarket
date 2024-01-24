package users

import (
	"stockmarket/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context, db *gorm.DB) []models.User {

	// get all users from gorm
	var users []models.User
	db.Find(&users)

	return users // passed into templates
}

func Show(c *gin.Context, db *gorm.DB) models.User {
	id := c.Param("id")

	var user models.User
	db.First(&user, id)

	return user // passed into templates

}