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

func (user *User) ActiveGamePlayer(db *gorm.DB) (uint, error) {
	var player Player
	err := db.
		Joins("Game").
		Where("user_id = ? AND active = ?", user.ID, true).
		First(&player).
		Error

	return player.ID, err
}

func (user *User) CreatePlayer(game Game, db *gorm.DB) (Player, error) {

	player := Player{
		Game:   game,
		User:   *user,
		Active: true,
		Cash:   100000,
	}
	err := db.Create(&player).Error

	if err != nil {
		fmt.Println("error creating player:", err)
		return Player{}, err
	}

	// create player stocks
	if err != nil {
		fmt.Println("error fetching stocks:", err)
		return Player{}, err
	}

	game_stocks := []GameStock{}

	err = db.Where("game_id = ?", game.ID).Find(&game_stocks).Error

	if err != nil {
		fmt.Println("error fetching game stocks:", err)
		return Player{}, err
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
			return Player{}, err
		}
	}

	return player, err
}

func (current_user *User) SetActiveGame(game Game, db *gorm.DB) (Player, error) {

	player, err := current_user.GetPlayer(game, db)

	// if gorm no record error
	if err != nil {
		fmt.Println("player does not exist, creating")
		player, err = current_user.CreatePlayer(game, db)
	}

	if err != nil {
		fmt.Println("error creating player:", err)
		return Player{}, err
	}

	if !player.Active {
		fmt.Println("setting active game for:", current_user.ID, game.ID)
		player.Active = true

		err = db.Model(&player).Where("id = ?", player.ID).Update("active", true).Error
	}

	if err != nil {
		fmt.Println("error setting active game for:", current_user.ID, game.ID, err)
		return Player{}, err
	}

	fmt.Println("unsetting active game for other games", current_user.ID, game.ID)
	err = db.Model(&Player{}).Where("user_id = ? AND game_id != ?", current_user.ID, game.ID).Update("active", false).Error

	if err != nil {
		fmt.Println("error unsetting active game for other games:", err)
		return Player{}, err
	}

	fmt.Println("set successfully")

	return player, nil
}

func GenerateSessionToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString, err
}
