package routes

import (
	"context"
	"stockmarket/models"
	templ "stockmarket/templates"
	"stockmarket/templates/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// get user from gorm db context with id
		var user models.User
		db.First(&user, id)

		pageComponent := users.Show(user)
		baseComponent := templ.Base("User - id", pageComponent)

		baseComponent.Render(context.Background(), c.Writer)
	})

	r.GET("/user-card/:id", func(c *gin.Context) {
		id := c.Param("id")

		// get user from gorm db context with id
		var user models.User
		db.First(&user, id)

		userComponent := users.Card(user)

		userComponent.Render(context.Background(), c.Writer)
	})

	return r
}
