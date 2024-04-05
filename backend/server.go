package main

import (
	"log"
	"os"
	"stockmarket/database"
	"stockmarket/router"
	"stockmarket/websockets"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func init() {
	LoadEnvVariables()
}

func main() {

	switch os.Getenv("ENVIRONMENT") {
	case "production":
		database.SetupProdDb()
	case "development":
		database.SetupDevDb()
	}

	websockets.InitializeGameHub()

	websockets.InitializeGameClosedHub()

	r := router.SetupRoutes()
	r.Static("/static", "./static")
	var port = os.Getenv("PORT")

	if port == "" {
		panic("port was not found")
	}

	r.Run(":" + port)

}
