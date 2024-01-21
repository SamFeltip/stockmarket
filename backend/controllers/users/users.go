package users

import (
	"stockmarket/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context, db *gorm.DB) []models.User {

	// get all users from gorm
	var users []models.User
	db.Find(&users)

	return users // passed into templates
}

func Show(c *gin.Context, db *gorm.DB) models.User {
	id := c.Param("id")

	var user models.User
	db.First(&user, id)

	return user // passed into templates

}

// func Create(c *gin.Context, db *gorm.DB) (models.User, error) {
// 	// name := c.PostForm("name")
// 	// profile := c.PostForm("profile")
// 	// password := c.PostForm("password")

// 	_, err = models.DoesUserExist(db, name)

// 	if !errors.Is(err, gorm.ErrRecordNotFound) {
// 		return models.User{}, errors.New("user already exists")
// 	}

// 	filePath := "/static/imgs/" + profile + "profile.png"

// 	user := models.User{
// 		Name:        name,
// 		Password:    password,
// 		ProfileRoot: filePath,
// 	}

// 	db.Create(&user)

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"ID":  user.ID,
// 		"exp": time.Now().Add(time.Hour).Unix(),
// 	})

// 	tokenString, err := token.SignedString([]byte("your-secret-key"))

// 	if err != nil {
// 		return models.User{}, err
// 	}

// 	c.SetCookie("token", tokenString, 3600, "", "", false, true)

// 	return user, nil // passed into templates
// }
