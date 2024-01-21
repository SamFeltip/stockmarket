package authorisation

import (
	"fmt"
	"net/http"
	"stockmarket/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// [POST] /signup
func Signup(c *gin.Context, db *gorm.DB) {

	var body struct {
		Name     string
		Profile  string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	filePath := "/static/imgs/" + body.Profile + "profile.png"

	fmt.Println(filePath)

	user := models.User{Name: body.Name, Password: string(hash), ProfileRoot: filePath}
	result := db.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create new user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": nil})

}

func Login(c *gin.Context, db *gorm.DB) {

	var body struct {
		Name     string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user, err := models.DoesUserExist(db, body.Name)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find user",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to compare password",
		})
		return
	}

	token, err := models.GenerateSessionToken(user)
	if err != nil {
		fmt.Println("Failed to generate session token:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate session token",
		})
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Authorisation", token, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"error": nil})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"current_user": user,
	})
}
