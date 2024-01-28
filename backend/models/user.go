package models

import (
	"fmt"
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

func (user *User) ActiveGame(db *gorm.DB) (Game, error) {
	var player Player
	fmt.Println("finding active game", user.ID)
	err := db.Preload("Game").Where("user_id = ? AND active = ?", user.ID, true).First(&player).Error

	return player.Game, err
}

func (current_user *User) SetActiveGame(game Game, db *gorm.DB) error {

	player, find_err := GetPlayer(&game, current_user, db)

	if find_err == gorm.ErrRecordNotFound {
		fmt.Println("could not find player, creating:", find_err)
		player = Player{
			Game:   game,
			User:   *current_user,
			Active: true,
		}
		create_err := db.Create(&player).Error

		if create_err != nil {
			fmt.Println("error creating player:", create_err)
			return find_err
		}
	} else if find_err != nil {
		fmt.Println("unexpected error fetching player:", find_err)
		return find_err
	}

	// set all players of user to inactive
	db.Model(&Player{}).Where("user_id = ?", current_user.ID).Update("active", false)

	player.Active = true
	db.Save(&player)

	return nil
}

func GenerateSessionToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString, err
}
