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
		c.Header("HX-Reconnect", "/login")
		c.Redirect(http.StatusFound, "/login")
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
		c.Header("HX-Reconnect", "/login")
		c.Redirect(http.StatusFound, "/login")
		return models.User{}, fmt.Errorf("expired cookie")
	}

	// find user with token sub
	var user models.User
	err = db.Where(claims["sub"]).First(&user).Error

	if err != nil {
		fmt.Println("could not find user: ", err)
		c.Header("HX-Reconnect", "/login")
		c.Redirect(http.StatusFound, "/login")
	}

	return user, nil
}

func AuthIsPlaying(c *gin.Context) {

	user, err := TestAuth(c)

	if err != nil {
		fmt.Println("could not find user: ", err)
		c.Header("HX-Reconnect", "/login")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Set("user", user)

	db := database.GetDb()
	playerID, err := user.ActiveGamePlayer(db)

	if err != nil {
		fmt.Println("user is not participating in a game (RequireAuth)", err)
		c.Header("HX-Reconnect", "/")
		c.Redirect(http.StatusFound, "/")
		return
	}

	player, err := models.FindPlayer(playerID, db)

	if err != nil {
		fmt.Println("error fetching player:", err)
		c.Header("HX-Reconnect", "/")
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.Set("player", player)

	game, err := models.FindGame(player.GameID, db)

	if err != nil {
		fmt.Println("error fetching game:", err)
		c.Header("HX-Reconnect", "/")
		c.Redirect(http.StatusFound, "/")
		return
	}

	fmt.Println("active game:", game.ID)
	c.Set("gameID", game.ID)

	c.Next()
}

func AuthIsLoggedIn(c *gin.Context) {

	user, err := TestAuth(c)

	if err != nil {
		fmt.Println("could not find user: ", err)
		c.Header("HX-Reconnect", "/login")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Set("user", user)

	c.Next()
}

func RequireAuthWebsocket(c *gin.Context) {
	db := database.GetDb()

	user, err := TestAuth(c)

	gameID := c.Param("gameID")

	if gameID == "" {
		fmt.Println("no gameID given in websocket connection request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no gameID"})
		return
	}

	game, err := models.LoadGame(gameID, db)

	if err != nil {
		fmt.Println("error fetching game:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "no game found"})
		return
	}

	player, err := user.SetActiveGame(gameID, db)

	if err != nil {
		fmt.Println("could not set active game for user: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Set("user", user)
	c.Set("player", player)
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
	fmt.Println("auth current player")

	db := database.GetDb()

	c.Request.ParseForm()
	user, err := TestAuth(c)

	if err != nil {
		fmt.Println("invalid credentials", err)
		c.Header("HX-Reconnect", "/login")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	fmt.Println("username:", user.Name)

	gameID := c.PostForm("gameID")

	if gameID == "" {
		fmt.Println("no gameID given in form which requires auth current player")
		c.Header("HX-Reconnect", "/login")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	game, err := models.LoadGame(gameID, db)

	if err != nil {
		fmt.Println("error fetching game:", err)
		return
	}

	c.Set("gameID", game.ID)

	if game.CurrentUser.ID != user.ID {
		fmt.Println("user is not current player", game.CurrentUserID, user.ID)
		return
	}

	c.Next()
}
