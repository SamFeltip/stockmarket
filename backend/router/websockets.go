package router

import (
	controllers "stockmarket/controllers/websockets"
	"stockmarket/middleware"

	"github.com/gin-gonic/gin"
)

func CreateWebsocketRoutes() {

	r.GET("/load-players",
		func(c *gin.Context) { middleware.RequireAuthWebsocket(c) },
		func(c *gin.Context) {
			httpResponseCode, response := controllers.LoadPlayers(c)
			c.JSON(httpResponseCode, response)
		})

	r.GET("/update-difficulty",
		func(c *gin.Context) { middleware.AuthCurrentPlayer(c) },
		func(c *gin.Context) {
			httpResponseCode, response := controllers.UpdateDifficulty(c)
			c.JSON(httpResponseCode, response)

		})
}
