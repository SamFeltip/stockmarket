package database

import (
	"fmt"
	"log"
	"stockmarket/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupDb() *gorm.DB {

	dsn := "host=localhost user=me password=def78-brglger-45y$u3g dbname=postgres port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.DB()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the PostgreSQL database!")

	// Perform migrations
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	// Perform migrations
	err = db.AutoMigrate(&models.Game{})
	if err != nil {
		log.Fatal(err)
	}

	// Perform migrations
	err = db.AutoMigrate(&models.Player{})
	if err != nil {
		log.Fatal(err)
	}

	return db

}

func GetDb() *gorm.DB {
	return db
}
