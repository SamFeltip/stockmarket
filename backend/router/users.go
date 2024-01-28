package router

import (
	"context"
	controllers "stockmarket/controllers/users"
	templates "stockmarket/templates/users"

	"github.com/gin-gonic/gin"
)

func CreateUserRoutes() {

	r.GET("/users", func(c *gin.Context) {
		users := controllers.Index(c)

		pageComponent := templates.Index(users)
		RenderWithTemplate(pageComponent, "Users", c)
	})

	r.GET("/users/show/:id", func(c *gin.Context) {
		user := controllers.Show(c)

		pageComponent := templates.Show(user)
		RenderWithTemplate(pageComponent, "User - id", c)

	})

	r.GET("/users/card/:id", func(c *gin.Context) {

		user := controllers.Show(c)
		userComponent := templates.Card(user)
		userComponent.Render(context.Background(), c.Writer)

	})
}
