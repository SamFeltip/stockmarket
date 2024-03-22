package games

import (
	"fmt"
	"stockmarket/database"
	"stockmarket/models"
	templates "stockmarket/templates/games"
	userTemplates "stockmarket/templates/users"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Show(db *gorm.DB, c *gin.Context) templ.Component {

	gameIDcontext, exists := c.Get("gameID")

	if !exists {
		fmt.Println("no game found")
		return templates.NoGame()
	}

	gameID := gameIDcontext.(string)

	cu, exists := c.Get("user")

	if !exists {
		fmt.Println("no user found")
		return templates.Error(fmt.Errorf("no user found"))
	}

	current_user := cu.(models.User)

	fmt.Println("run setActiveGame", current_user.Name)
	current_player, err := current_user.SetActiveGame(gameID, db)

	if err != nil {
		fmt.Println("error setting active game:", err)
		return templates.Error(err)
	}

	fmt.Println("player fetched:", current_player.User.Name, current_player.Game.ID)
	fmt.Println("player active:", current_player.Active)

	players, err := models.GetPlayers(gameID, db)
	userCardList := userTemplates.CardListSocket(players)

	err = BroadcastUpdatePlayersList(gameID, userCardList)

	if err != nil {
		fmt.Println("error broadcasting new player:", err)
		return templates.Error(err)
	}

	fmt.Println("player stocks" + strconv.Itoa(len(current_player.PlayerStocks)))

	game, err := models.FindGame(gameID, db)

	if game.Status == string(models.Playing) {

		gameDisplay, err := models.LoadGameDisplay(gameID, db)

		if err != nil {
			fmt.Println("error fetching game:", err)
			gameWrapper := templates.NoGame()
			return gameWrapper
		}

		currentPlayerDisplay, err := models.LoadCurrentPlayerDisplay(current_player.ID, db)

		pageComponent := templates.Playing(gameDisplay, currentPlayerDisplay)
		return pageComponent
	}

	if game.Status == string(models.Closed) {

		gameInsights, err := models.GetGameInsights(game.ID, db)

		if err != nil {
			fmt.Println("error getting game insights:", err)
			return templates.Error(err)
		}

		pageComponent := templates.Closed(gameInsights, game.GameStocks, game.Players)
		return pageComponent
	}

	pageComponent := templates.Waiting(game)
	return pageComponent
}

func Create(c *gin.Context) (models.Game, error) {
	db := database.GetDb()

	code := strings.ToLower(c.PostForm("code"))
	periodCountStr := c.PostForm("difficulty")

	if code == "" || periodCountStr == "" {
		fmt.Println("no code or periodCount in form")
		return models.Game{}, fmt.Errorf("no code or periodCount in form")
	}

	periodCount, err := strconv.Atoi(periodCountStr)
	if err != nil {
		// handle error, e.g. return an error response
		fmt.Println("couldnt convert to int")
		return models.Game{}, err
	}

	cu, _ := c.Get("user")
	current_user := cu.(models.User)

	if err != nil {
		fmt.Println("error creating game stocks, creating empty set:", err)
		return models.Game{}, err
	}

	fmt.Println("create game:", code, periodCount)

	game := models.Game{
		ID:          code,
		PeriodCount: periodCount,
		Status:      string(models.Waiting),
		CurrentUser: current_user,
	}

	db.Create(&game)

	// get all stocks
	fmt.Println("create game stocks: ", code)
	game_stocks, err := models.CreateGameStocks(code, db)

	if err != nil {
		fmt.Println("error creating game stocks:", err)
		return models.Game{}, err
	}

	game.GameStocks = game_stocks
	db.Save(&game)

	if err != nil {
		fmt.Println("error setting active game:", err)
		return models.Game{}, err
	}

	_, err = models.NewFeedItem(game, 0, models.StartGame, models.Player{}, models.GameStock{}, db)

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
