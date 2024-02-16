package router

import (
	controllers "stockmarket/controllers/insights"
	templates "stockmarket/templates/insights"

	"github.com/gin-gonic/gin"
)

func CreateInsightRoutes() {

	r.GET("/insights",
		func(c *gin.Context) {
			insights := controllers.Index(c)

			insightTemplate := templates.Index(insights)
			RenderWithTemplate(insightTemplate, "Insights", c)
		})
}
