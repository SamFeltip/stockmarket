package games

import (
	"fmt"
	"stockmarket/database"
	"stockmarket/models"
	templates "stockmarket/templates/games"

	"github.com/a-h/templ"
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

	game, err := models.FindGame(gameID, db)

	if err != nil {
		fmt.Println("could not find game", err)
		return nil, err
	}

	err = db.Model(models.Game{}).Where("lower(id) = lower(?)", gameID).Update("status", string(models.Playing)).Error

	if err != nil {
		fmt.Println("could not update game status")
		return nil, err
	}

	fmt.Println("game status updated:", gameID)

	players, err := models.GetPlayers(gameID, db)

	err = game.GeneratePlayerInsights(players, db)

	if err != nil {
		fmt.Println("could not generate player insights", err)
		return nil, err
	}

	fmt.Println("player insights made:", game.ID)
	err = BroadcastUpdatePlayBoard(game.ID)

	if err != nil {
		fmt.Println("could not broadcast start play")
		return nil, err
	}

	baseComponent := templates.WaitingLoading()

	return baseComponent, nil
}

func PlayAction(gameID string, current_user models.User, db *gorm.DB) (templ.Component, error) {

	game, err := models.FindGame(gameID, db)

	if err != nil {
		fmt.Println("could not find game", err)
		return templates.Error(err), err
	}

	_, err = models.NewFeedItemMessage(game.ID, game.CurrentPeriod, models.PlayerPass, current_user, db)

	if err != nil {
		fmt.Println("could not create new feed item", err)
		return templates.Error(err), err
	}

	template, err := CheckForMarketClose(game.ID, db)

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
	_, err = models.UpdateCurrentUser(game.ID, db)

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

func NextPeriod(gameID string, db *gorm.DB) (templ.Component, error) {

	game, err := models.FindGame(gameID, db)

	if err != nil {
		fmt.Println("could not find game", err)
		return templates.Error(err), err
	}

	err = game.UpdatePeriod(db)

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
