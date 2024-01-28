package middleware

import (
	"fmt"
	"net/http"
	"os"
	"stockmarket/database"
	"stockmarket/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuth(c *gin.Context) (models.User, error) {

	db := database.GetDb()

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

func RequireAuth(c *gin.Context) {

	user, err := TestAuth(c)

	if err != nil {
		fmt.Println("could not find user: ", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Set("user", user)

	c.Next()
}

func RequireAuthWebsocket(c *gin.Context) {
	db := database.GetDb()

	user, err := TestAuth(c)

	if err != nil {
		fmt.Println("could not find user: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	game, err := user.ActiveGame(db)
	if err != nil {
		fmt.Println("user is not participating in a game", err)
		// write a http response and return
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not participating in a game"})
		return
	}

	c.Set("user", user)
	c.Set("game", game)

	c.Next()
}

func SoftAuth(c *gin.Context) {

	user, err := TestAuth(c)

	if err != nil {
		fmt.Println("could not find user: ", err)
	}

	c.Set("user", user)
	c.Next()

}

func AuthCurrentPlayer(c *gin.Context) {
	db := database.GetDb()

	c.Request.ParseForm()
	user, err := TestAuth(c)
	if err != nil {
		fmt.Println("invalid credentials")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "it is not your turn",
		})
		return
	}

	gameID := c.Request.Form["gameID"][0]

	var game models.Game
	err = db.Where("lower(id) = lower(?) AND current_user_id = ?", gameID, user.ID).First(&game).Error

	if err != nil {
		fmt.Println("could not find game")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "it is not your turn",
		})
		return
	}

	c.Next()
}
