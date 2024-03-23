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
	current_user_id, err := models.UpdateCurrentUser(game.ID, db)
	if err != nil {
		t.Fatalf("UpdateCurrentUser failed: %v", err)
	}

	// Check that the current user was updated correctly
	assert.Equal(t, user2.ID, current_user_id)

	// Retrieve the game from the database
	var retrievedGame models.Game
	db.First(&retrievedGame, "id = ?", game.ID)

	// Check that the current user was updated correctly
	assert.Equal(t, user2.ID, retrievedGame.CurrentUserID)

	// Call UpdateCurrentUser again to cycle the current user
	current_user_id, err = models.UpdateCurrentUser(game.ID, db)
	if err != nil {
		t.Fatalf("UpdateCurrentUser failed: %v", err)
	}

	// Check that the current user was cycled correctly
	assert.Equal(t, user1.ID, current_user_id)

	// Retrieve the game from the database again
	db.First(&retrievedGame, "id = ?", game.ID)

	// Check that the current user was cycled correctly
	assert.Equal(t, user1.ID, retrievedGame.CurrentUserID)

}

/*

func (game *Game) UpdatePeriod(db *gorm.DB) error {

	type GameStockChange struct {
		TotalChange float64
		GameStockID uint
		Value       float64
	}

	var gameStockChanges []GameStockChange

	err := db.Table("player_stocks as ps").
		Select("sum(i.value) as total_change, gs.ID as game_stock_id, gs.value").
		Joins("left join player_insights as pi on pi.player_stock_id = ps.id").
		Joins("left join insights as i on i.id = pi.insight_id").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Joins("inner join stocks as s on s.id = gs.stock_id").
		Where("gs.game_id = ?", "some").
		Group("s.name, s.image_path, gs.value, s.variation").
		Order("s.variation").
		Scan(&gameStockChanges).Error

	if err != nil {
		fmt.Println("could not get game stock changes", err)
		return err
	}

	// loop through gameStockChanges and update gameStocks
	for _, gameStockChange := range gameStockChanges {
		gameStock := GameStock{}
		err = db.
			Model(&gameStock).
			Where("id = ?", gameStockChange.GameStockID).
			Update("value", gameStockChange.Value+gameStockChange.TotalChange).Error

		if err != nil {
			fmt.Println("could not update game stock", err)
			return err
		}
	}

	err = game.GeneratePlayerInsights(db)

	if err != nil {
		fmt.Println("could not generate player insights", err)
		return err
	}

	game.CurrentPeriod++
	game.Status = string(Playing)

	err = db.Save(&game).Error
	return err
}

*/

func TestUpdatePeriod(t *testing.T) {
	db := database.SetupTestDb(logger.Info)
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

	// Call UpdatePeriod to update the period
	err = game.UpdatePeriod(db)
	if err != nil {
		t.Fatalf("UpdatePeriod failed: %v", err)
	}

	// Retrieve the game from the database
	var retrievedGame models.Game
	db.First(&retrievedGame, "id = ?", "gameTest")

	// Check that the period was updated correctly
	assert.Equal(t, 1, retrievedGame.PeriodCount)
	assert.Equal(t, "playing", retrievedGame.Status)

	// asset that the game stock values have been updated
	var gameStocks []models.GameStock
	db.Find(&gameStocks, "game_id = ?", "gameTest")

	for _, gameStock := range gameStocks {
		assert.NotEqual(t, 100, gameStock.Value)
	}

}
