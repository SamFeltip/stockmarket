package main

import (
	"stockmarket/database"
	"stockmarket/router"
)

func main() {

	db := database.SetupDb()
	r := router.SetupRoutes(db)
	r.Static("/static", "./static")
	r.Run(":4040")

}
