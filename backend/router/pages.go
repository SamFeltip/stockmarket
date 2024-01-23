package router

import (
	"stockmarket/middleware"
	templates "stockmarket/templates/pages"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePageRoutes(db *gorm.DB, r *gin.Engine) {

	r.GET("/",
		func(c *gin.Context) { middleware.RequireAuth(c, db) },
		func(c *gin.Context) {

			pageComponent := templates.Index()
			RenderWithTemplate(pageComponent, "Stockmarket", c)
		})

}
