package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stockmarket/database"
	"stockmarket/models"
	"time"

	templates "stockmarket/templates/games"

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

func AuthIsPlaying(c *gin.Context) {

	user, err := TestAuth(c)

	if err != nil {
		fmt.Println("could not find user: ", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Set("user", user)

	db := database.GetDb()
	player, err := user.ActiveGamePlayer(db)
	game := player.Game

	fmt.Println("active game:", game.ID)

	if err != nil {
		fmt.Println("user is not participating in a game (RequireAuth)", err)
		// write a http response and return
		game = models.Game{}
		return
	}
	c.Set("game", game)
	c.Set("player", player)

	c.Next()
}

func AuthIsLoggedIn(c *gin.Context) {

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
	fmt.Println("username for active game:", user.Name)
	player, err := user.ActiveGamePlayer(db)

	if err != nil {
		fmt.Println("user is not participating in a game (RequireAuthWebsocket)", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user found in request context"})
		return
		// game = models.Game{}
	}

	c.Set("user", user)
	c.Set("player", player)
	c.Set("game", player.Game)

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
	fmt.Println("auth current player")

	db := database.GetDb()

	c.Request.ParseForm()
	user, err := TestAuth(c)

	if err != nil {
		fmt.Println("invalid credentials", err)

		pageComponent := templates.Error(err)
		ctx := context.Background()
		pageComponent.Render(ctx, c.Writer)
		return
	}

	fmt.Println("username:", user.Name)

	gameID := c.PostForm("gameID")

	if gameID == "" {
		fmt.Println("no gameID given in form which requires auth current player")
		pageComponent := templates.Error(fmt.Errorf("no gameID"))
		ctx := context.Background()
		pageComponent.Render(ctx, c.Writer)
		return
	}

	var game models.Game
	err = db.Where("lower(games.id) = lower(?) AND current_user_id = ?", gameID, user.ID).First(&game).Error

	if err != nil {
		fmt.Println("could not find game", err)
		pageComponent := templates.Error(err)
		ctx := context.Background()
		pageComponent.Render(ctx, c.Writer)
		return
	}

	c.Next()
}
