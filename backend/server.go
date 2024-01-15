package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Define a GORM model
type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func readFileAsString(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func main() {

	dsn := "host=localhost user=me password=def78-brglger-45y$u3g dbname=postgres port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	fmt.Println("Connected to the PostgreSQL database!")

	// Perform migrations
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")

	// html := template.Must(template.ParseFiles("templates/base/index.html"))
	// r.SetHTMLTemplate(html)

	r.GET("/hello", func(c *gin.Context) {

		_, err := template.ParseFiles("templates/base/index.html", "templates/pages/index.html")

		if err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "layout", gin.H{
			"Title": "Homes",
		})
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.HTML(http.StatusOK, "users/index", gin.H{
			"Title": "User Page",
			"ID":    id,
		})
	})

	r.GET("/get-info", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/name_card", gin.H{
			"Name":  "Sam",
			"Phone": "+447565328118",
			"Email": "sf@gmail.com",
		})
	})

	r.Run(":4040")

	// e.Logger.Fatal(e.Start(":4040"))

}
