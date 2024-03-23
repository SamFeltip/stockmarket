package models

import (
	"fmt"

	"gorm.io/gorm"
)

type FeedItem struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`

	Message   string
	Title     string
	ImageRoot string
	Colour    string

	PlayerStock   PlayerStock
	PlayerStockID uint `gorm:"default:null"`

	Game   Game
	GameID string `gorm:"foreignKey:ID"`

	Period int
}

func LoadFeedItems(game_id string, db *gorm.DB) ([]FeedItem, error) {

	var feed_items []FeedItem

	err := db.
		Preload("Player.User").
		Preload("Player").
		Where("game_id = ?", game_id).
		Order("created_at desc").
		Find(&feed_items).Error

	if err != nil {
		fmt.Println("error loading feed items:", err)
		return nil, err
	}

	return feed_items, nil
}

type FeedItemMessage string

var StartGame FeedItemMessage = "startGame"
var PlayerPlay FeedItemMessage = "playerPlay"
var PlayerPass FeedItemMessage = "playerPass"
var PeriodNew FeedItemMessage = "periodNew"

func NewFeedItemMessage(gameID string, currentPeriod int, feedItemMessage FeedItemMessage, user User, db *gorm.DB) (FeedItem, error) {

	feed_item := FeedItem{
		GameID: gameID,
		Period: currentPeriod,
	}

	if feedItemMessage == StartGame {
		feed_item.Message = "game started"
		feed_item.Title = "game started"
		feed_item.ImageRoot = "/static/imgs/icons/Handshake.svg"
		feed_item.Colour = "grey"
	}

	if feedItemMessage == PlayerPass {
		feed_item.Message = "passed their go"
		feed_item.Title = user.Name
		feed_item.ImageRoot = user.ProfileRoot
	}

	fmt.Println("creating new FeedItem", feed_item.Message)
	err := db.Create(&feed_item).Error

	if err != nil {
		return FeedItem{}, err
	}

	return feed_item, nil
}

func NewFeedItem(quantity int, playerStockID uint, db *gorm.DB) (FeedItem, error) {

	type FeedItemData struct {
		UserName        string
		UserProfileRoot string
		StockName       string
		GameID          string
		CurrentPeriod   int
	}

	var feed_item_data FeedItemData

	err := db.Table("player_stocks as ps").
		Select("u.name as user_name, u.profile_root as user_profile_root, s.Name as stock_name, gs.game_id, g.current_period").
		Joins("inner join game_stocks as gs on gs.id = ps.game_stock_id").
		Joins("inner join games as g on g.id = gs.game_id").
		Joins("inner join stocks as s on s.id = gs.stock_id").
		Joins("inner join players as p on p.id = ps.player_id").
		Joins("inner join users as u on u.id = p.user_id").
		Where("ps.ID = ?", playerStockID).
		First(&feed_item_data).Error

	if err != nil {
		fmt.Println("error loading feed item data:", err)
		return FeedItem{}, err
	}

	feed_item := FeedItem{
		PlayerStockID: playerStockID,
		Period:        feed_item_data.CurrentPeriod,

		Title:     feed_item_data.UserName,
		ImageRoot: feed_item_data.UserProfileRoot,
	}

	if quantity > 0 {
		feed_item.Message = fmt.Sprintf("bought %d shares in %s", quantity, feed_item_data.StockName)
	} else if quantity < 0 {
		feed_item.Message = fmt.Sprintf("sold %d shares in %s", quantity*-1, feed_item_data.StockName)
	}

	fmt.Println("creating new FeedItem", feed_item.Message)
	err = db.Create(&feed_item).Error

	if err != nil {
		return FeedItem{}, err
	}

	return feed_item, nil
}
