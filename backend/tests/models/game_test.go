package models

import (
	"stockmarket/database"
	"stockmarket/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestCreateGame(t *testing.T) {
	db := database.SetupTestDb(logger.Silent)
	defer database.UndoMigrations(db)

	// Create a user and save it to the database
	user := models.User{
		Name:     "user1",
		Password: "password",
	}

	db.Create(&user)

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

	db.Create(&stock1)
	db.Create(&stock2)

	// Call CreateGame to create a game
	game, err := models.CreateGame("gameTest", 1, user, db)

	if err != nil {
		t.Fatalf("CreateGame failed: %v", err)
	}

	// Check that the game was created correctly
	assert.Equal(t, 1, game.PeriodCount)
	assert.Equal(t, user.ID, game.CurrentUserID)
	assert.Equal(t, "waiting", game.Status)

	// Retrieve the game from the database
	var retrievedGame models.Game
	db.First(&retrievedGame, "id = ?", "gameTest")

	// Check that a game was created with this code
	assert.Equal(t, "gameTest", retrievedGame.ID)

}

func TestUpdateCurrentUser(t *testing.T) {
	db := database.SetupTestDb(logger.Info)
	defer database.UndoMigrations(db)

	// Create two users and save them to the database
	user1 := models.User{
		Name:        "user1",
		Password:    "password1",
		ProfileRoot: "profile1",
	}
	user2 := models.User{
		Name:        "user2",
		Password:    "password2",
		ProfileRoot: "profile2",
	}

	db.Create(&user1)
	db.Create(&user2)

	game, _ := models.CreateGame("game1", 1, user1, db)

	player1 := models.Player{
		User:   user1,
		Game:   game,
		Active: true,
		Cash:   100000,
	}

	player2 := models.Player{
		User:   user2,
		Game:   game,
		Active: true,
		Cash:   100000,
	}

	db.Create(&player1)
	db.Create(&player2)

	game.Players = append(game.Players, player1)
	game.Players = append(game.Players, player2)

	// Call UpdateCurrentUser to update the user
	err := game.UpdateCurrentUser(db)
	if err != nil {
		t.Fatalf("UpdateCurrentUser failed: %v", err)
	}

	// Retrieve the game from the database
	var retrievedGame models.Game
	db.First(&retrievedGame, "id = ?", game.ID)

	// Check that the current user was updated correctly
	assert.Equal(t, user2.ID, retrievedGame.CurrentUserID)

	// Call UpdateCurrentUser again to cycle the current user
	err = game.UpdateCurrentUser(db)
	if err != nil {
		t.Fatalf("UpdateCurrentUser failed: %v", err)
	}

	// Retrieve the game from the database again
	db.First(&retrievedGame, "id = ?", game.ID)

	// Check that the current user was cycled correctly
	assert.Equal(t, user1.ID, retrievedGame.CurrentUserID)

}
