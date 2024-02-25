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

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BroadcastUpdatePlayersList(players []models.Player) error {

	userCardList := userTemplates.CardListSocket(players)

	buffer := &bytes.Buffer{}
	userCardList.Render(context.Background(), buffer)

	latestPlayer := players[len(players)-1]

	broadcastMessage := websocketModels.BroadcastMessage{
		Game:   latestPlayer.Game,
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

func BroadcastStartPlay(game models.Game) error {

	fmt.Println("broadcasting start play: capturing playing socket template")

	broadcastMessage := websocketModels.BroadcastMessage{
		Game:    game,
		Buffer:  nil,
		Message: "start play",
	}

	fmt.Println("broadcasting start play: sending playing socket template")
	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel
	return nil
}

func Show(c *gin.Context) templ.Component {
	fmt.Println("show!!!!")
	db := database.GetDb()

	gameID := c.Param("id")

	game, err := models.GetGame(gameID, db)

	if err != nil {
		fmt.Println("error fetching game:", err)
		return templates.NoGame()
	}

	fmt.Println("show game:", game.ID)
	c.Set("game", game)

	cu, _ := c.Get("user")
	current_user := cu.(models.User)

	fmt.Println("setting active game", current_user.Name)
	player, err := current_user.SetActiveGame(game, db)

	if err != nil {
		fmt.Println("error setting active game:", err)
		return templates.Error(err)
	}

	c.Set("player", player)

	fmt.Println("player fetched:", player.User.Name, player.Game.ID)

	game.Players = append(game.Players, player)

	err = BroadcastUpdatePlayersList(game.Players)

	if err != nil {
		fmt.Println("error broadcasting new player:", err)
		return templates.Error(err)
	}

	return templates.IngamePage(player, game)
}

func Create(c *gin.Context) (models.Game, error) {
	db := database.GetDb()

	code := c.PostForm("code")
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

	return game, nil // passed into templates
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
