package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"stockmarket/models"
	templ "stockmarket/templates"
	pages "stockmarket/templates/pages"
	users "stockmarket/templates/users"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		user := models.User{
			Name: "Sam",
		}

		pageComponent := pages.Greeting(user)
		baseComponent := templ.Base("root", pageComponent)

		baseComponent.Render(context.Background(), c.Writer)
	})

	r.GET("/hello", func(c *gin.Context) {

		pageComponent := pages.Welcome()
		baseComponent := templ.Base("welcome user!", pageComponent)

		baseComponent.Render(context.Background(), c.Writer)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// get user from gorm db context with id
		var user models.User
		db.First(&user, id)

		pageComponent := users.Show(user)
		baseComponent := templ.Base("User - id", pageComponent)

		baseComponent.Render(context.Background(), c.Writer)
	})

	r.GET("/user-card/:id", func(c *gin.Context) {
		id := c.Param("id")

		// get user from gorm db context with id
		var user models.User
		db.First(&user, id)

		userComponent := users.Card(user)

		userComponent.Render(context.Background(), c.Writer)

	})

	r.Run(":4040")

	// e.Logger.Fatal(e.Start(":4040"))

}
