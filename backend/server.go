package main

import (
	"fmt"
	"log"
	"net/http"
	"stockmarket/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
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

func main() {
	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Group level middleware
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "joe" && password == "secret" {
			return true, nil
		}
		return false, nil
	}))

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

	// e.Static("/", "static")

	// Little bit of middlewares for housekeeping
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))

	e.Static("/static", "static")

	//This will initiate our template renderer
	template.NewTemplateRenderer(e, "pages/*.html")

	e.GET("/hello", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/get-info", func(c echo.Context) error {
		res := map[string]interface{}{
			"Name":  "Sam",
			"Phone": "+447565328118",
			"Email": "sf@gmail.com",
		}
		return c.Render(http.StatusOK, "name_card", res)
	})

	e.Logger.Fatal(e.Start(":4040"))

}
