package router

import (
	controller "stockmarket/controllers/authorisation"
	"stockmarket/middleware"
	templates "stockmarket/templates/authorisation"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateAuthRoutes(db *gorm.DB, r *gin.Engine) {

	r.GET("/signup", func(c *gin.Context) {

		pageComponent := templates.Signup()
		RenderWithTemplate(pageComponent, "Signup", c)

	})

	r.GET("/login", func(c *gin.Context) {

		pageComponent := templates.Login()
		RenderWithTemplate(pageComponent, "Login", c)

	})

	r.POST("/signup", func(c *gin.Context) { controller.Signup(c, db) })

	r.POST("/login", func(c *gin.Context) { controller.Login(c, db, controller.SignupBody{}) })

	r.GET(
		"/validate",
		func(c *gin.Context) { middleware.RequireAuth(c, db) },
		controller.Validate,
	)

	r.POST(
		"/logout",
		controller.Logout,
	)
}
