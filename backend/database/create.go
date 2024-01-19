package database

import (
	"fmt"
	"log"
	"stockmarket/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	return db

}
