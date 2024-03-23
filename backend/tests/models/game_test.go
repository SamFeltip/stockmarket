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

func TestGeneratePlayerInsights(t *testing.T) {
	db := database.SetupTestDb(logger.Warn)
	defer database.UndoMigrations(db)

	// Create a create_game_user and save it to the database
	create_game_user := models.User{
		Name:     "user1",
		Password: "password",
	}

	another_user := models.User{
		Name:     "user2",
		Password: "password",
	}

	db.Create(&create_game_user)
	db.Create(&another_user)

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

	insights := []models.Insight{
		{StockID: stock1.ID, Description: "insight1", Value: 1},
		{StockID: stock1.ID, Description: "insight2", Value: 2},
		{StockID: stock2.ID, Description: "insight3", Value: 1.5},
		{StockID: stock2.ID, Description: "insight4", Value: 2.5},
		{StockID: stock1.ID, Description: "insight5", Value: 0.5},
		{StockID: stock2.ID, Description: "insight6", Value: 2},
		{StockID: stock1.ID, Description: "insight7", Value: -1.5},
		{StockID: stock2.ID, Description: "insight8", Value: 1},
		{StockID: stock1.ID, Description: "insight9", Value: -2.5},
		{StockID: stock2.ID, Description: "insight10", Value: 3},
		{StockID: stock1.ID, Description: "insight11", Value: 1.5},
		{StockID: stock2.ID, Description: "insight12", Value: 0.5},
		{StockID: stock1.ID, Description: "insight13", Value: 2.5},
		{StockID: stock2.ID, Description: "insight14", Value: -1},
		{StockID: stock1.ID, Description: "insight15", Value: -0.5},
		{StockID: stock2.ID, Description: "insight16", Value: -2},
		{StockID: stock1.ID, Description: "insight17", Value: -1.5},
		{StockID: stock2.ID, Description: "insight18", Value: -0.5},
		{StockID: stock1.ID, Description: "insight19", Value: 1.5},
		{StockID: stock2.ID, Description: "insight20", Value: 0.5},
		{StockID: stock1.ID, Description: "insight21", Value: 2.5},
		{StockID: stock2.ID, Description: "insight22", Value: -1},
		{StockID: stock1.ID, Description: "insight23", Value: -0.5},
		{StockID: stock2.ID, Description: "insight24", Value: -2},
		{StockID: stock1.ID, Description: "insight25", Value: -1.5},
		{StockID: stock2.ID, Description: "insight26", Value: -0.5},
		{StockID: stock1.ID, Description: "insight27", Value: 1.5},
		{StockID: stock2.ID, Description: "insight28", Value: 0.5},
		{StockID: stock1.ID, Description: "insight29", Value: 2.5},
		{StockID: stock2.ID, Description: "insight30", Value: -1},
	}

	db.Create(&insights)

	game := models.Game{
		ID:            "gameTest",
		PeriodCount:   1,
		CurrentUserID: create_game_user.ID,
		Status:        "playing",
	}

	db.Create(&game)

	player1 := models.Player{
		UserID: create_game_user.ID,
		GameID: game.ID,
		Active: true,
		Cash:   100000,
	}

	player2 := models.Player{
		UserID: another_user.ID,
		GameID: game.ID,
		Active: true,
		Cash:   100000,
	}

	db.Create(&player1)
	db.Create(&player2)

	gameStock1 := models.GameStock{
		GameID:  game.ID,
		StockID: stock1.ID,
		Value:   100,
	}

	gameStock2 := models.GameStock{
		GameID:  game.ID,
		StockID: stock2.ID,
		Value:   200,
	}

	db.Create(&gameStock1)
	db.Create(&gameStock2)

	playerStock1 := models.PlayerStock{
		PlayerID:    player1.ID,
		GameStockID: gameStock1.ID,
		Quantity:    10,
	}

	playerStock2 := models.PlayerStock{
		PlayerID:    player1.ID,
		GameStockID: gameStock2.ID,
		Quantity:    1,
	}

	playerStock3 := models.PlayerStock{
		PlayerID:    player2.ID,
		GameStockID: gameStock1.ID,
		Quantity:    10,
	}

	playerStock4 := models.PlayerStock{
		PlayerID:    player2.ID,
		GameStockID: gameStock2.ID,
		Quantity:    1,
	}

	db.Create(&playerStock1)
	db.Create(&playerStock2)
	db.Create(&playerStock3)
	db.Create(&playerStock4)

	err := game.GeneratePlayerInsights(
		[]models.Player{player1, player2}, db)

	if err != nil {
		t.Fatalf("GeneratePlayerInsights failed: %v", err)
	}

	// Retrieve the player insights from the database
	// retrieve pi where pi.game_stock_id = game_stock_id where game_stock.game_id = game.id
	var playerInsights []models.PlayerInsight
	err = db.Table("player_insights").
		Joins("inner join player_stocks on player_insights.player_stock_id = player_stocks.id").
		Joins("inner join game_stocks on player_stocks.game_stock_id = game_stocks.id").
		Where("game_stocks.game_id = ?", game.ID).
		Find(&playerInsights).Error

	if err != nil {
		t.Fatalf("could not retrieve player insights: %v", err)
	}

	// Check that the player insights were generated correctly
	assert.Equal(t, 20, len(playerInsights))
}
