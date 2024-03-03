package models

import (
	"stockmarket/database"
	"stockmarket/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestCreatePlayer(t *testing.T) {
	db := database.SetupTestDb(logger.Info)
	defer database.UndoMigrations(db)

	stock1 := models.Stock{
		Name:          "stock1",
		StartingValue: 100,
		ImagePath:     "image1",
		Variation:     0.5,
	}

	stock2 := models.Stock{
		Name:          "stock2",
		StartingValue: 200,
		ImagePath:     "image2",
		Variation:     0.5,
	}

	game_create_user := models.User{
		Name:     "user1",
		Password: "password",
	}
	db.Create(&game_create_user)

	game := models.Game{
		ID:          "game1",
		Difficulty:  1,
		Status:      string(models.Waiting),
		CurrentUser: game_create_user,
	}
	db.Create(&game)

	db.Create(&stock1)
	db.Create(&stock2)

	gameStock1 := models.GameStock{
		Game:  game,
		Stock: stock1,
		Value: stock1.StartingValue,
	}

	gameStock2 := models.GameStock{
		Game:  game,
		Stock: stock2,
		Value: stock2.StartingValue,
	}

	db.Create(&gameStock1)
	db.Create(&gameStock2)

	game_create_user_player := models.Player{
		User:   game_create_user,
		Game:   game,
		Active: true,
		Cash:   100000,
	}

	db.Create(&game_create_user_player)

	new_user := models.User{
		Name:     "user2",
		Password: "password",
	}

	db.Create(&new_user)

	// Call CreatePlayer to create a player
	new_player, err := new_user.CreatePlayer(&game, db)

	if err != nil {
		t.Fatalf("CreatePlayer failed: %v", err)
	}

	// Check that the player was created correctly
	if new_player == nil {
		t.Fatalf("CreatePlayer failed: %v", game_create_user_player)
	}

	// check player is in game
	game_player_ref, _ := game.GetPlayer(&new_user)

	assert.Equal(t, game_player_ref.ID, new_player.ID)

	// check player is created in db
	var retrievedPlayer models.Player
	db.First(&retrievedPlayer, "id = ?", new_player.ID)

	assert.Equal(t, new_player.ID, retrievedPlayer.ID)

	// check player stocks are in db
	var playerStocks []models.PlayerStock
	db.Where("player_id = ?", new_player.ID).Find(&playerStocks)

	assert.Equal(t, 2, len(playerStocks))

}