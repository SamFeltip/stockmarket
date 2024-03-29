package games

import (
	"fmt"
	"stockmarket/database"
	"stockmarket/models"
	templates "stockmarket/templates/games"
	userTemplates "stockmarket/templates/users"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Show(gameID string, current_user models.User, db *gorm.DB) templ.Component {

	fmt.Println("run setActiveGame", current_user.Name)
	current_player, err := current_user.SetActiveGame(gameID, db)

	if err != nil {
		fmt.Println("error setting active game:", err)
		return templates.Error(err)
	}

	fmt.Println("player fetched:", current_player.ID)
	fmt.Println("player active:", current_player.Active)

	players, err := models.LoadPlayerDisplays(gameID, db)

	if err != nil {
		fmt.Println("error getting players:", err)
		return templates.Error(err)
	}

	userCardList := userTemplates.CardListSocket(players)

	err = BroadcastUpdatePlayersList(gameID, userCardList)

	if err != nil {
		fmt.Println("error broadcasting new player:", err)
		return templates.Error(err)
	}

	fmt.Println("player stocks", strconv.Itoa(len(current_player.PlayerStocks)))

	game, err := models.FindGame(gameID, db)

	if err != nil {
		fmt.Println("error finding game:", err)
		return templates.Error(err)
	}

	switch game.Status {
	case string(models.Playing):
		gameDisplay, err := models.LoadGameDisplay(gameID, db)

		if err != nil {
			fmt.Println("error fetching game:", err)
			gameWrapper := templates.NoGame()
			return gameWrapper
		}

		currentPlayerDisplay, err := models.LoadCurrentPlayerDisplay(current_player.ID, db)

		if err != nil {
			fmt.Println("error fetching current player:", err)
			gameWrapper := templates.NoGame()
			return gameWrapper
		}

		pageComponent := templates.Playing(gameDisplay, currentPlayerDisplay, players)
		return pageComponent
	case string(models.Closed):
		gameInsights, err := models.GetGameInsights(game.ID, db)

		if err != nil {
			fmt.Println("error getting game insights:", err)
			return templates.Error(err)
		}

		gameStockDisplays, err := models.LoadGameStockDisplays(gameID, db)

		if err != nil {
			fmt.Println("error loading game stock displays:", err)
			return templates.Error(err)
		}

		playerDisplays, err := models.LoadPlayerDisplays(gameID, db)

		if err != nil {
			fmt.Println("error loading player displays:", err)
			return templates.Error(err)
		}

		pageComponent := templates.Closed(gameID, gameInsights, gameStockDisplays, playerDisplays)
		return pageComponent
	default:
		fmt.Println("players length", len(players))
		pageComponent := templates.Waiting(game, players, current_user.ID)
		return pageComponent
	}
}

func Create(code string, periodCount int, current_user models.User, db *gorm.DB) (models.Game, error) {

	fmt.Println("create game:", code, periodCount)

	game := models.Game{
		ID:            code,
		PeriodCount:   periodCount,
		Status:        string(models.Waiting),
		CurrentUserID: current_user.ID,
	}

	db.Create(&game)

	// get all stocks
	fmt.Println("create game stocks: ", code)
	_, err := models.CreateGameStocks(code, db)

	if err != nil {
		fmt.Println("error creating game stocks:", err)
		return models.Game{}, err
	}

	db.Save(&game)

	if err != nil {
		fmt.Println("error setting active game:", err)
		return models.Game{}, err
	}

	_, err = models.NewFeedItemMessage(game.ID, game.CurrentPeriod, models.StartGame, models.User{}, db)

	if err != nil {
		fmt.Println("could not create new feed item", err)
		return models.Game{}, err
	}

	return game, nil
}

func New(c *gin.Context, db *gorm.DB) {
}

func Index(c *gin.Context) templ.Component {
	db := database.GetDb()

	// get all games from gorm
	var games []models.Game
	db.Find(&games)

	pageComponent := templates.Index(games)

	return pageComponent // passed into templates
}
