package users

import (
	"fmt"
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

func New(c *gin.Context, db *gorm.DB) models.User {
	name := c.PostForm("name")
	profile := c.PostForm("profile")
	password := c.PostForm("password")

	// print profile
	fmt.Println(profile)

	filePath := "/static/imgs/" + profile + "profile.png"

	user := models.User{
		Name:        name,
		Password:    password,
		ProfileRoot: filePath,
	}

	db.Create(&user)

	return user // passed into templates
}
