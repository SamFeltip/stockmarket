package router

import (
	"context"
	controllers "stockmarket/controllers/users"
	templates "stockmarket/templates/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUserRoutes(db *gorm.DB, r *gin.Engine) {

	r.GET("/users", func(c *gin.Context) {
		users := controllers.Index(c, db)

		pageComponent := templates.Index(users)
		RenderWithTemplate(pageComponent, "Users", c)
	})

	r.GET("/users/show/:id", func(c *gin.Context) {
		user := controllers.Show(c, db)

		pageComponent := templates.Show(user)
		RenderWithTemplate(pageComponent, "User - id", c)

	})

	r.GET("/users/card/:id", func(c *gin.Context) {

		user := controllers.Show(c, db)
		userComponent := templates.Card(user)
		userComponent.Render(context.Background(), c.Writer)

	})

	// r.POST("/users/new", func(c *gin.Context) {

	// 	user, err := controllers.Create(c, db)

	// 	if err != nil {
	// 		fmt.Println("error creating user:", err)
	// 		userComponent := templates.Card(user)
	// 		RenderWithTemplate(userComponent, "error!", c)
	// 	}

	// 	userComponent := templates.Card(user)
	// 	RenderWithTemplate(userComponent, "User - id", c)

	// })
}
