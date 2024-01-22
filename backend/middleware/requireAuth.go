package middleware

import (
	"fmt"
	"net/http"
	"os"
	"stockmarket/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func TestAuth(c *gin.Context, db *gorm.DB) (models.User, error) {

	// Get the cookie of req
	tokenString, err := c.Cookie("Authorisation")

	if err != nil {
		fmt.Println("unauthorised: ", err)
		return models.User{}, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println("invalid token: ", err)
		return models.User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		fmt.Println("invalid token: ")

		return models.User{}, fmt.Errorf("invalid token")
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		fmt.Println("expired cookie: ")
		c.Redirect(http.StatusFound, "/login")
		return models.User{}, fmt.Errorf("expired cookie")
	}

	// find user with token sub
	var user models.User
	err = db.Where(claims["sub"]).First(&user).Error

	if err != nil {
		fmt.Println("could not find user: ", err)
		return models.User{}, err
	}

	return user, nil
}

func RequireAuth(c *gin.Context, db *gorm.DB) {

	user, err := TestAuth(c, db)

	if err != nil {
		fmt.Println("could not find user: ", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Set("user", user)

	c.Next()
}

func SoftAuth(c *gin.Context, db *gorm.DB) {

	user, err := TestAuth(c, db)

	if err != nil {
		fmt.Println("could not find user: ", err)
	}

	c.Set("user", user)
	c.Next()

}
