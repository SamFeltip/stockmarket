package router

import (
	controllers "stockmarket/controllers/websockets"
	"stockmarket/middleware"
	"stockmarket/websockets"

	"github.com/gin-gonic/gin"
)

func CreateWebsocketRoutes() {
	websockets.InitializeHub()

	r.GET("/ws",
		func(c *gin.Context) { middleware.RequireAuthWebsocket(c) },
		func(c *gin.Context) { controllers.ServeWs(c) })
}
