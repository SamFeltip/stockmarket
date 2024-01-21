package main

import (
	"log"
	"os"
	"stockmarket/database"
	"stockmarket/router"

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

	db := database.SetupDb()
	r := router.SetupRoutes(db)
	r.Static("/static", "./static")
	var port = os.Getenv("PORT")

	if port == "" {
		panic("port was not found")
	}

	r.Run(":" + port)

}
