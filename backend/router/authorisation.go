package router

import (
	controller "stockmarket/controllers/authorisation"
	templates "stockmarket/templates/authorisation"

	"github.com/gin-gonic/gin"
)

func CreateAuthRoutes() {

	r.GET("/signup", func(c *gin.Context) {

		pageComponent := templates.Signup()
		RenderWithTemplate(pageComponent, "Signup", c)

	})

	r.GET("/login", func(c *gin.Context) {

		pageComponent := templates.Login()
		RenderWithTemplate(pageComponent, "Login", c)

	})

	r.POST("/signup", func(c *gin.Context) { controller.Signup(c) })

	r.POST("/login", func(c *gin.Context) { controller.Login(c, controller.SignupBody{}) })

	r.GET(
		"/validate",
		controller.Validate,
	)

	r.GET(
		"/logout",
		func(c *gin.Context) {
			c.SetCookie("Authorisation", "", -1, "", "", false, true)

			pageComponent := templates.Logout()

			RenderWithTemplate(pageComponent, "Logout", c)

		})

}
