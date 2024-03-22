package games

import (
	"fmt"
	"stockmarket/database"
	"stockmarket/models"
	templates "stockmarket/templates/games"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateGamePeriodCount(gameID string, periodCount int) (templ.Component, error) {

	db := database.GetDb()

	game := models.Game{}
	err := db.Model(&game).Where("lower(games.id) = lower(?)", gameID).Update("period_count", periodCount).Error

	errMsg := ""

	if err != nil {
		fmt.Println("could not update game PeriodCount")
		return nil, err
	}

	err = db.Model(&game).Where("lower(games.id) = lower(?)", gameID).First(&game).Error

	if err != nil {
		fmt.Printf("Error reloading game: %v", err)
		errMsg = "Error reloading game"
		return nil, err
	}

	err = BroadcastUpdatePeriodCount(game)

	if err != nil {
		fmt.Println("could not broadcast PeriodCount update")
		errMsg = "could not broadcast PeriodCount update"
		return nil, err
	}

	baseComponent := templates.PeriodCountOptions(game, errMsg)
	return baseComponent, nil
}

func StartGame(gameID string) (templ.Component, error) {

	db := database.GetDb()

	err := db.Model(models.Game{}).Where("lower(id) = lower(?)", gameID).Update("status", string(models.Playing)).Error

	if err != nil {
		fmt.Println("could not update game status")
		return nil, err
	}

	fmt.Println("game status updated:", gameID)

	game, err := models.LoadGame(gameID, db)

	if err != nil {
		fmt.Println("could not find game", err)
		return nil, err
	}

	game.GeneratePlayerInsights(db)

	fmt.Println("player insights made:", game.ID)
	err = BroadcastUpdatePlayBoard(game.ID)

	if err != nil {
		fmt.Println("could not broadcast start play")
		return nil, err
	}

	baseComponent := templates.WaitingLoading()

	return baseComponent, nil
}

func PlayAction(c *gin.Context, db *gorm.DB) (templ.Component, error) {
	gameAction, _ := c.GetPostForm("game_action")

	gameIDcontext, exists := c.Get("gameID")

	if !exists {
		fmt.Println("game doesn't exist in context")
		return templates.Error(fmt.Errorf("game doesn't exist in context")), fmt.Errorf("game doesn't exist in context")
	}

	gameID := gameIDcontext.(string)
	game, err := models.LoadGameDisplay(gameID, db)

	if err != nil {
		fmt.Println("error fetching game:", err)
		return templates.NoGame(), err
	}

	cu, exists := c.Get("user")

	if !exists {
		fmt.Println("no user found")
		return templates.Error(fmt.Errorf("no user found")), fmt.Errorf("no user found")
	}

	current_user := cu.(models.User)

	player, err := current_user.GetPlayer(gameID, db)

	_, err = models.NewFeedItem(game, 0, models.PlayerPass, current_user, player.ID, models.GameStock{}, db)

	if err != nil {
		fmt.Println("could not create new feed item", err)
		return templates.Error(err), err
	}

	if err != nil {
		fmt.Println("could not create new feed item", err)
		return templates.Error(err), err
	}

	template, err := CheckForMarketClose(game, db)

	if err != nil {
		fmt.Println("could not check for market close", err)
		return templates.Error(err), err
	}

	if template != nil {
		fmt.Println("market closed")
		return template, nil
	}

	fmt.Println("market not closed")

	// update current user
	err = game.UpdateCurrentUser(db)

	if err != nil {
		fmt.Println("could not update current player", err)
		return templates.Error(err), err
	}

	err = BroadcastUpdatePlayBoard(game.ID)

	if err != nil {
		fmt.Println("could not broadcast update board", err)
		return templates.Error(err), err
	}

	return templates.Loading(), nil
}

func NextPeriod(c *gin.Context, db *gorm.DB) (templ.Component, error) {
	cg, exists := c.Get("game")

	if !exists {
		fmt.Println("game doesn't exist in context")
		return templates.Error(fmt.Errorf("game doesn't exist in context")), fmt.Errorf("game doesn't exist in context")
	}

	game := cg.(models.Game)

	if game.Status != string(models.Closed) {
		fmt.Println("game not closed")
		return templates.Error(fmt.Errorf("game not closed")), fmt.Errorf("game not closed")
	}

	err := game.UpdatePeriod(db)

	if err != nil {
		fmt.Println("could not update period", err)
		return templates.Error(err), err
	}

	err = BroadcastUpdatePlayBoard(game.ID)

	if err != nil {
		fmt.Println("could not broadcast period update", err)
		return templates.Error(err), err
	}

	return templates.Loading(), nil
}
