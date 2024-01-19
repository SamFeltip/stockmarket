package router

import (
	"stockmarket/models"
	templates "stockmarket/templates/pages"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePageRoutes(db *gorm.DB, r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		user := models.User{
			Name: "Sam",
		}

		pageComponent := templates.Greeting(user)
		RenderWithTemplate(pageComponent, "Stockmarket", c)
	})

	r.GET("/signup", func(c *gin.Context) {

		pageComponent := templates.Signup()
		RenderWithTemplate(pageComponent, "Signup", c)

	})
}
