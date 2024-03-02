package games

import (
	"bytes"
	"context"
	"fmt"
	"stockmarket/database"
	"stockmarket/models"
	websocketModels "stockmarket/models/websockets"
	templates "stockmarket/templates/games"
	userTemplates "stockmarket/templates/users"
	"stockmarket/websockets"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BroadcastUpdatePlayersList(game *models.Game) error {

	userCardList := userTemplates.CardListSocket(game.Players)

	buffer := &bytes.Buffer{}
	userCardList.Render(context.Background(), buffer)

	if len(game.Players) == 0 {
		fmt.Println("no players to broadcast")
		return nil
	}

	broadcastMessage := websocketModels.BroadcastMessage{
		Game:   *game,
		Buffer: buffer,
	}

	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel

	return nil
}

func BroadcastUpdateDifficulty(game models.Game) error {

	difficultyDisplay := templates.DifficultyOptionsSocket(game)

	buffer := &bytes.Buffer{}
	difficultyDisplay.Render(context.Background(), buffer)

	broadcastMessage := websocketModels.BroadcastMessage{
		Game:   game,
		Buffer: buffer,
	}

	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel

	return nil
}

func BroadcastUpdateBoard(game models.Game) error {

	fmt.Println("broadcasting show board: capturing playing socket template")

	broadcastMessage := websocketModels.BroadcastMessage{
		Game:    game,
		Buffer:  nil,
		Message: "game board",
	}

	fmt.Println("broadcasting show board: sending playing socket template")
	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel
	return nil
}

func Show(db *gorm.DB, c *gin.Context) templ.Component {

	cg, exists := c.Get("game")
	game := cg.(models.Game)

	if !exists {
		fmt.Println("no game found")
		return templates.NoGame()
	}

	cu, _ := c.Get("user")
	current_user := cu.(models.User)

	fmt.Println("run setActiveGame", current_user.Name)
	current_player, err := game.SetActiveGame(&current_user, db)

	if err != nil {
		fmt.Println("error setting active game:", err)
		return templates.Error(err)
	}

	c.Set("player", current_player)

	fmt.Println("player fetched:", current_player.User.Name, current_player.Game.ID)
	fmt.Println("player active:", current_player.Active)

	// game.Players = append(game.Players, player)

	err = BroadcastUpdatePlayersList(&game)

	if err != nil {
		fmt.Println("error broadcasting new player:", err)
		return templates.Error(err)
	}

	fmt.Println("player stocks" + strconv.Itoa(len(current_player.PlayerStocks)))

	return templates.IngamePage(game, current_player)
}

func Create(c *gin.Context) (models.Game, error) {
	db := database.GetDb()

	code := strings.ToLower(c.PostForm("code"))
	difficultyStr := c.PostForm("difficulty")

	if code == "" || difficultyStr == "" {
		fmt.Println("no code or difficulty in form")
		return models.Game{}, fmt.Errorf("no code or difficulty in form")
	}

	difficulty, err := strconv.Atoi(difficultyStr)
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

	fmt.Println("create game:", code, difficulty)

	game := models.Game{
		ID:          code,
		Difficulty:  difficulty,
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
