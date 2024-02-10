package router

import (
	"fmt"
	controllers "stockmarket/controllers/users"
	templates "stockmarket/templates/users"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func CreateUserRoutes() {

	r.GET("/users", func(c *gin.Context) {
		users := controllers.Index(c)

		pageComponent := templates.Index(users)
		RenderWithTemplate(pageComponent, "Users", c)
	})

	r.GET("users/show/:id", func(c *gin.Context) {
		user, err := controllers.Show(c)

		var pageComponent templ.Component

		if err != nil {
			fmt.Println("user could not be found")
			pageComponent = templates.NoUser()
		} else {
			pageComponent = templates.Show(user)
		}

		RenderWithTemplate(pageComponent, "User - show", c)
	})

}
