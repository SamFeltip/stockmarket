package models

import (
	"stockmarket/database"
	"stockmarket/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestCreatePlayer(t *testing.T) {
	db := database.SetupTestDb(logger.Warn)
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
		PeriodCount: 1,
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
	new_player, err := new_user.CreatePlayer(game.ID, db)

	if err != nil {
		t.Fatalf("CreatePlayer failed: %v", err)
	}

	// Check that the player was created correctly
	if new_player.ID == 0 {
		t.Fatalf("CreatePlayer failed: %v", game_create_user_player)
	}

	// check player is created in db
	var retrievedPlayer models.Player
	db.First(&retrievedPlayer, "id = ?", new_player.ID)

	assert.Equal(t, new_player.ID, retrievedPlayer.ID)

	// check player stocks are in db
	var playerStocks []models.PlayerStock
	db.Where("player_id = ?", new_player.ID).Find(&playerStocks)

	assert.Equal(t, 2, len(playerStocks))

}

// test LoadCurrentPlayerDisplay

func TestLoadCurrentPlayerDisplay(t *testing.T) {
	db := database.SetupTestDb(logger.Warn)

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
		PeriodCount: 1,
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

	playerDisplay, err := models.LoadCurrentPlayerDisplay(game_create_user_player.ID, db)

	if err != nil {
		t.Fatalf("LoadCurrentPlayerDisplay failed: %v", err)
	}

	assert.Equal(t, game_create_user_player.ID, playerDisplay.UserID)

	if len(playerDisplay.PlayerStocks) != 2 {
		t.Fatalf("LoadCurrentPlayerDisplay failed: %v", playerDisplay.PlayerStocks)
	}

}

func TestLoadPlayerDisplays(t *testing.T) {

	db := database.SetupTestDb(logger.Warn)

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

	db.Create(&stock1)
	db.Create(&stock2)

	game_create_user := models.User{
		Name:     "user1",
		Password: "password",
	}

	db.Create(&game_create_user)

	game := models.Game{
		ID:          "game1",
		PeriodCount: 1,
		Status:      string(models.Waiting),
		CurrentUser: game_create_user,
	}

	db.Create(&game)

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

	playerStock1 := models.PlayerStock{
		Player:    game_create_user_player,
		GameStock: gameStock1,
		Quantity:  1,
	}

	playerStock2 := models.PlayerStock{
		Player:    game_create_user_player,
		GameStock: gameStock2,
		Quantity:  1,
	}

	db.Create(&playerStock1)
	db.Create(&playerStock2)

	player_display, err := models.LoadCurrentPlayerDisplay(game_create_user_player.ID, db)

	if err != nil {
		t.Fatalf("LoadCurrentPlayerDisplay failed: %v", err)
	}

	assert.Equal(t, 2, len(player_display.PlayerStocks))
}
