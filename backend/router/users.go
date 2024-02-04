package router

import (
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

}
