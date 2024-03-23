package database

import (
	"fmt"
	"log"
	"os"
	"stockmarket/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func UndoMigrations(db *gorm.DB) {
	db.Migrator().DropTable(&models.PlayerInsight{})
	db.Migrator().DropTable(&models.Insight{})
	db.Migrator().DropTable(&models.PlayerStock{})
	db.Migrator().DropTable(&models.GameStock{})
	db.Migrator().DropTable(&models.Stock{})
	db.Migrator().DropTable(&models.Player{})
	db.Migrator().DropTable(&models.Game{})
	db.Migrator().DropTable(&models.User{})
}

func SetupTestDb(log_mode logger.LogLevel) *gorm.DB {
	dsn := "host=localhost user=me password=password dbname=stockmarket_test port=5433 sslmode=disable"
	return SetupDb(dsn, log_mode)
}

func SetupDevDb() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	return SetupDb(dsn, logger.Info)
}

func SetupProdDb() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	return SetupDb(dsn, logger.Error)
}

func SetupDb(dsn string, log_mode logger.LogLevel) *gorm.DB {

	newdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(log_mode),
	})

	db = newdb

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

	// Perform migrations
	err = db.AutoMigrate(&models.Stock{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.GameStock{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.PlayerStock{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.Insight{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.PlayerInsight{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.FeedItem{})
	if err != nil {
		log.Fatal(err)
	}

	return db

}

func GetDb() *gorm.DB {
	return db
}
