package router

import (
	"context"
	"fmt"
	controllers "stockmarket/controllers/player_stocks"
	"stockmarket/database"
	"stockmarket/middleware"
	models "stockmarket/models"
	gameTemplates "stockmarket/templates/games"
	templates "stockmarket/templates/player_stocks"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePlayerStockRoutes() {

	r.GET("/player_stocks/show/:playerStockID",
		func(c *gin.Context) { middleware.AuthIsLoggedIn(c) },
		func(c *gin.Context) {
			db := database.GetDb()

			playerStockIDString := c.Param("playerStockID")

			// convert playerStockIDString to uint
			playerStockIDuint64, err := strconv.ParseUint(playerStockIDString, 10, 64)

			if err != nil {
				fmt.Println("error converting playerStockIDString to uint", err)
				pageComponent := gameTemplates.Error(fmt.Errorf("error converting playerStockIDString to uint"))
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			playerStockID := uint(playerStockIDuint64)

			var playerStockPlayer models.PlayerStockPlayerResult

			fmt.Println("getting player info from player stock:", playerStockID)
			// player info
			db.Table("player_stocks as ps").
				Select("ps.quantity as stocks_held, (gs.value * ps.quantity) as stock_value, p.cash").
				Joins("inner join game_stocks as gs on ps.game_stock_id = gs.id").
				Joins("inner join players as p on p.id = ps.player_id").
				Where("ps.id = ?", playerStockIDString).
				Scan(&playerStockPlayer)

			var playerStockDisplay models.PlayerStockDisplay

			fmt.Println("getting total insights for player stock")
			// total insights for player stock
			err = db.Table("player_stocks as ps").
				Select("ps.ID, gs.game_id, sum(i.value) as total_insight, gs.value as game_stock_value, gs.game_id, s.name as stock_name, s.image_path as stock_image_path").
				Joins("left join player_insights as pi on pi.player_stock_id = ps.id").
				Joins("left join insights as i on i.id = pi.insight_id").
				Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
				Joins("inner join stocks as s on s.id = gs.stock_id").
				Where("ps.id = ?", playerStockIDString).
				Group("ps.id, gs.value, s.name, s.image_path, gs.game_id").
				Scan(&playerStockDisplay).Error

			if err != nil {
				fmt.Println("error fetching player stock preview", err)
				pageComponent := gameTemplates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			var investors []models.InvestorResult

			// all investors
			db.Table("player_stocks as ps").
				Select("u.name, u.profile_root, ps.quantity").
				Joins("inner join player_stocks as psl on ps.game_stock_id = psl.game_stock_id").
				Joins("inner join players as p on p.id = ps.player_id").
				Joins("inner join users as u on u.id = p.user_id").
				Where("psl.id = ?", playerStockIDString).
				Scan(&investors)

			var insightResults []models.InsightResult

			// my insights
			db.Table("player_insights as pi").
				Select("i.description, i.value").
				Joins("inner join player_stocks as ps on ps.id = pi.player_stock_id").
				Joins("inner join insights as i on pi.insight_id = i.id").
				Where("ps.id = ?", playerStockIDString).
				Scan(&insightResults)

			var stockInfoResult models.StockInfoResult

			// stock info
			db.Table("player_stocks as ps").
				Select("(100000 - sum(ps.quantity)) as shares_available, s.variation").
				Joins("inner join player_stocks as psl on ps.game_stock_id = psl.game_stock_id").
				Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
				Joins("inner join stocks as s on s.id = gs.stock_id").
				Where("psl.id = ?", playerStockIDString).
				Group("s.variation").
				Scan(&stockInfoResult)

			// is current player
			var result struct {
				IsCurrentPlayer bool
			}

			db.Table("player_stocks as ps").
				Select("g.current_user_id = u.id as is_current_player").
				Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
				Joins("inner join games as g on g.id = gs.game_id").
				Joins("inner join players as p on p.id = ps.player_id").
				Joins("inner join users as u on p.user_id = u.id").
				Where("ps.id = ?", playerStockIDString).
				Scan(&result)

			isCurrentPlayer := result.IsCurrentPlayer

			pageComponent := templates.Show(
				playerStockID,
				playerStockPlayer,
				playerStockDisplay,
				investors,
				insightResults,
				stockInfoResult,
				isCurrentPlayer)

			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)

		})

	r.POST("/player_stocks/edit",
		func(c *gin.Context) { middleware.AuthCurrentPlayer(c) },
		func(c *gin.Context) {
			db := database.GetDb()
			playerStockIDString := c.PostForm("PlayerStockID")
			playerStockQuantityAdd := c.PostForm("PlayerStockQuantityAdd")

			if playerStockIDString == "" || playerStockQuantityAdd == "" {
				fmt.Println("no playerStockID or playerStockQuantityAdd in form")
				pageComponent := gameTemplates.Error(fmt.Errorf("no playerStockID or playerStockQuantityAdd"))
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			fmt.Println("form data gathered: playerStockID:", playerStockIDString, "playerStockQuantityAdd:", playerStockQuantityAdd)

			gameID := c.PostForm("gameID")
			mode := c.PostForm("mode")

			playerStockID64, err := strconv.ParseUint(playerStockIDString, 10, 32)

			if err != nil {
				fmt.Println("could not parse id", err)
				pageComponent := gameTemplates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			playerStockID := uint(playerStockID64)

			// parse QuantityAdd to int and add to player stock . quantity
			quantityAdd, err := strconv.Atoi(playerStockQuantityAdd)
			if err != nil {
				fmt.Println("could not parse new quantity to int", err)
				pageComponent := gameTemplates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			multiplier, err := strconv.Atoi(mode)
			if err != nil {
				fmt.Println("could not parse mode to int", err)
				pageComponent := gameTemplates.Error(err)
				ctx := context.Background()
				pageComponent.Render(ctx, c.Writer)
				return
			}

			pageComponent, err := controllers.Edit(playerStockID, gameID, quantityAdd, multiplier, db)
			ctx := context.Background()
			pageComponent.Render(ctx, c.Writer)

			if err != nil {
				fmt.Println("error editing player stock, don't broadcast", err)
				return
			}

		})
}
