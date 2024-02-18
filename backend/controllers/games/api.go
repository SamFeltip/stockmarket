package games

import (
	"fmt"
	"stockmarket/database"
	"stockmarket/models"
	templates "stockmarket/templates/games"

	"github.com/a-h/templ"
)

func UpdateGameDifficulty(gameID string, difficulty int) (templ.Component, error) {

	db := database.GetDb()

	game := models.Game{}
	err := db.Model(&game).Where("lower(games.id) = lower(?)", gameID).Update("difficulty", difficulty).Error

	errMsg := ""

	if err != nil {
		fmt.Println("could not update game difficulty")
		return nil, err
	}

	err = db.Model(&game).Where("lower(games.id) = lower(?)", gameID).First(&game).Error

	if err != nil {
		fmt.Printf("Error reloading game: %v", err)
		errMsg = "Error reloading game"
		return nil, err
	}

	err = BroadcastUpdateDifficulty(game)

	if err != nil {
		fmt.Println("could not broadcast difficulty update")
		errMsg = "could not broadcast difficulty update"
		return nil, err
	}

	baseComponent := templates.DifficultyOptions(game, errMsg)
	return baseComponent, nil
}

func StartGame(gameID string) (templ.Component, error) {

	db := database.GetDb()

	game := models.Game{}
	err := db.Model(&game).Where("id = lower(?)", gameID).Update("status", string(models.Playing)).Error

	if err != nil {
		fmt.Println("could not update game status")
		return nil, err
	}

	game, err = models.GetGame(gameID, db)

	if err != nil {
		fmt.Println("could not find game")
		return nil, err
	}

	game.GenerateInsights(db)

	fmt.Println("game updated:", game.ID)
	err = BroadcastStartPlay(game)

	if err != nil {
		fmt.Println("could not broadcast start play")
		return nil, err
	}

	baseComponent := templates.WaitingLoading()

	return baseComponent, nil
}
