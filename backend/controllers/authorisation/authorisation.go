package authorisation

import (
	"fmt"
	"net/http"
	"stockmarket/database"
	"stockmarket/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignupBody struct {
	Name     string `form:"Name"`
	Password string `form:"Password"`
	Profile  string `form:"Profile"`
}

// [POST] /signup
func Signup(c *gin.Context) {

	db := database.GetDb()
	var body SignupBody

	if c.Bind(&body) != nil {
		fmt.Println("Failed to read body on signup")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Name == "" || body.Password == "" {
		fmt.Println("missing name or password on signup")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing name or password",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		fmt.Println("Failed to hash password: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	filePath := "/static/imgs/" + body.Profile + "profile.png"

	user := models.User{Name: body.Name, Password: string(hash), ProfileRoot: filePath}
	err = db.Create(&user).Error

	if err != nil {
		fmt.Println("Failed to find user", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "username is taken",
		})
		return
	}

	fmt.Println("user signed up successfully")
	Login(c, body)

}

func Login(c *gin.Context, signupBody SignupBody) {
	db := database.GetDb()

	var loginBody struct {
		Name     string `form:"Name"`
		Password string `form:"Password"`
	}
	var err error

	if signupBody.Name != "" {
		loginBody.Name = signupBody.Name
		loginBody.Password = signupBody.Password
	} else {
		err = c.Bind(&loginBody)
	}

	if err != nil {
		fmt.Println("Failed to read body on login")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user, err := models.DoesUserExist(db, loginBody.Name)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find user",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginBody.Password))
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

	fmt.Println("user logged in successfully")
	c.JSON(http.StatusOK, gin.H{"error": nil})
}

func Validate(c *gin.Context) {
	cu, _ := c.Get("user")

	if cu == nil {
		fmt.Println("validate: no user found in context")
		cu = models.User{}
	}

	user := cu.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"current_user": user,
	})
}
