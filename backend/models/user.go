package models

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Define a GORM model
type User struct {
	gorm.Model
	ID          uint
	Name        string `gorm:"not null;unique"`
	Password    string `gorm:"not null"`
	ProfileRoot string `gorm:"not null"`
}

func DoesUserExist(db *gorm.DB, username string) (User, error) {

	var user User
	err := db.Where("lower(name) = lower(?)", username).First(&user).Error

	return user, err
}

func GenerateSessionToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString, err
}
