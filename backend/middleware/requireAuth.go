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

func RequireAuth(c *gin.Context, db *gorm.DB) {

	// Get the cookie of req
	tokenString, err := c.Cookie("Authorisation")

	if err != nil {
		fmt.Println("unauthorised: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no auth cookie",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	fmt.Println(claims)
	fmt.Println(ok)

	if !ok || !token.Valid {
		fmt.Println("invalid token: ", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		fmt.Println("expired cookie: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "expired session",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// find user with token sub
	var user models.User
	err = db.Where(claims["sub"]).First(&user).Error

	if err != nil {
		fmt.Println("could not find user: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user session",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)

	c.Next()
}
