package models

import (
	"stockmarket/database"
	"stockmarket/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestCreateGameStocks(t *testing.T) {

	db := database.SetupTestDb(logger.Silent)
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

	current_user := models.User{
		Name:     "user1",
		Password: "password",
	}

	game := models.Game{
		ID:          "game1",
		PeriodCount: 1,
		Status:      "waiting",
		CurrentUser: current_user,
	}

	db.Create(&game)

	game_stocks, _ := models.CreateGameStocks("game1", db)

	if len(game_stocks) != 2 {
		t.Fatalf("CreateGameStocks failed: %v", game_stocks)
	}

	// Check that the game stocks were created correctly
	assert.Equal(t, "game1", game_stocks[0].GameID)
	assert.Equal(t, stock1.ID, game_stocks[0].StockID)
	assert.Equal(t, "game1", game_stocks[1].GameID)
	assert.Equal(t, stock2.ID, game_stocks[1].StockID)

}
