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

// get the given user's active player
func (user *User) ActiveGamePlayer(db *gorm.DB) (uint, error) {
	var player Player
	err := db.
		Joins("Game").
		Where("user_id = ? AND active = ?", user.ID, true).
		First(&player).
		Error

	return player.ID, err
}

func (user *User) GetPlayer(gameID string, db *gorm.DB) (Player, error) {
	var player Player
	err := db.
		Where("user_id = ? AND game_id = ?", user.ID, gameID).
		First(&player).
		Error

	return player, err
}

/*
- create a new player object and connected player stocks

- adds player object to game object
*/
func (user *User) CreatePlayer(gameID string, db *gorm.DB) (Player, error) {

	player := Player{
		GameID: gameID,
		User:   *user,
		Active: true,
		Cash:   100000,
	}
	err := db.Create(&player).Error

	if err != nil {
		fmt.Println("error creating player:", err)
		return Player{}, err
	}

	game_stocks := []GameStock{}

	err = db.Where("game_id = ?", gameID).Find(&game_stocks).Error

	if err != nil {
		fmt.Println("error fetching game stocks:", err)
		return Player{}, err
	}

	for _, game_stock := range game_stocks {
		player_stock := PlayerStock{
			Player:    player,
			GameStock: game_stock,
			Quantity:  0,
		}

		err = db.Create(&player_stock).Error

		if err != nil {
			fmt.Println("error creating player stock:", err)
			return Player{}, err
		}
	}

	if err != nil {
		fmt.Println("error loading player:", err)
		return Player{}, err
	}

	return player, err
}

func GenerateSessionToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString, err
}
