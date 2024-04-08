package router

import (
	controllers "stockmarket/controllers/websockets"
	"stockmarket/middleware"

	"github.com/gin-gonic/gin"
)

func CreateWebsocketRoutes() {

	r.GET("/connected-game/:gameID",
		func(c *gin.Context) { middleware.RequireAuthWebsocket(c) },
		func(c *gin.Context) {
			httpResponseCode, response := controllers.ServeWGames(c)
			c.JSON(httpResponseCode, response)
		})

	r.GET("/hold-stock-waiting/:gameID",
		func(c *gin.Context) { middleware.RequireAuthWebsocket(c) },
		func(c *gin.Context) {
			httpResponseCode, response := controllers.ServeWInsights(c)
			c.JSON(httpResponseCode, response)
		})
}
