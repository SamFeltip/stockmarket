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
	err := db.Preload("Game").Preload("Game.CurrentUser").Where("user_id = ? AND active = ?", user.ID, true).First(&player).Error

	return player.Game, err
}

func (user *User) createPlayer(game Game, db *gorm.DB) (*Player, error) {
	player := Player{
		Game:   game,
		User:   *user,
		Active: true,
	}
	err := db.Create(&player).Error

	if err != nil {
		fmt.Println("error creating player:", err)
		return nil, err
	}

	// create player stocks
	if err != nil {
		fmt.Println("error fetching stocks:", err)
		return nil, err
	}

	for _, game_stock := range game.GameStocks {
		player_stock := PlayerStock{
			Player:    player,
			GameStock: game_stock,
			Quantity:  0,
		}

		err = db.Create(&player_stock).Error

		if err != nil {
			fmt.Println("error creating player stock:", err)
			return nil, err
		}
	}

	err = db.Model(&player).Preload("User").Preload("Game").First(&player).Error

	return &player, err
}

func (current_user *User) SetActiveGame(game Game, db *gorm.DB) error {

	player, find_err := GetPlayer(&game, current_user, db)

	if find_err == gorm.ErrRecordNotFound {
		fmt.Println("could not find player, creating:", find_err)
		new_player, create_err := current_user.createPlayer(game, db)

		if create_err != nil {
			fmt.Println("error creating player:", create_err)
			return create_err
		}

		player = *new_player

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
