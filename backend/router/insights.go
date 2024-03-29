package router

import (
	"context"
	"stockmarket/database"
	"stockmarket/models"
	templates "stockmarket/templates/insights"

	"github.com/gin-gonic/gin"
)

func CreateInsightRoutes() {

	r.GET("/insights",
		func(c *gin.Context) {
			db := database.GetDb()

			// get all games from gorm
			var insights []models.Insight
			db.Preload("Stock").Order("insights.ID asc").Find(&insights)

			insightTemplate := templates.Index(insights)
			RenderWithTemplate(insightTemplate, "Insights", c)
		})

	r.POST("/insights",
		func(c *gin.Context) {

			db := database.GetDb()

			// get insightID from post form
			insightID := c.PostForm("insightID")

			// get Description from post form
			description := c.PostForm("Description")

			if description == "" || insightID == "" {
				// return insights error
				pageComponent := templates.Message("invalid input")
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return

			}

			// update insight with gorm

			var insight models.Insight
			db.First(&insight, insightID)

			insight.Description = description
			err := db.Save(&insight).Error

			if err != nil {
				// return insights error
				pageComponent := templates.Message("error updating insight")
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			pageComponent := templates.Message("changed successfully")
			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)

		})
}
