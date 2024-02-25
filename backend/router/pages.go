package router

import (
	"stockmarket/middleware"
	templates "stockmarket/templates/pages"

	"github.com/gin-gonic/gin"
)

func CreatePageRoutes() {

	r.GET("/",
		func(c *gin.Context) { middleware.AuthIsPlaying(c) },
		func(c *gin.Context) {

			pageComponent := templates.Index()
			RenderWithTemplate(pageComponent, "Stockmarket", c)
		})

}
