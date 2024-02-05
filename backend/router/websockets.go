package router

import (
	controllers "stockmarket/controllers/websockets"
	"stockmarket/middleware"

	"github.com/gin-gonic/gin"
)

func CreateWebsocketRoutes() {

	r.GET("/connected-game",
		func(c *gin.Context) { middleware.RequireAuthWebsocket(c) },
		func(c *gin.Context) {
			httpResponseCode, response := controllers.ServeWs(c)
			c.JSON(httpResponseCode, response)
		})
}
