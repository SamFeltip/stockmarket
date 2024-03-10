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

	GameStock   GameStock
	GameStockID uint `gorm:"default:null"`

	Player   Player
	PlayerID uint `gorm:"default:null"`

	Game   Game
	GameID string

	Period int
}

func LoadFeedItems(game_id string, db *gorm.DB) ([]FeedItem, error) {

	var feed_items []FeedItem

	err := db.
		Preload("Player.User").
		Preload("GameStock").
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

func NewFeedItem(game Game, quantity int, feedItemMessage FeedItemMessage, player Player, game_stock GameStock, db *gorm.DB) (FeedItem, error) {

	feed_item := FeedItem{
		GameStock: game_stock,
		Player:    player,
		Game:      game,
		Period:    game.CurrentPeriod,
	}

	if quantity > 0 {
		feed_item.Message = fmt.Sprintf("bought %d shares in %s", quantity, game_stock.Stock.Name)
		feed_item.Title = player.User.Name
		feed_item.ImageRoot = player.User.ProfileRoot
	} else if quantity < 0 {
		feed_item.Message = fmt.Sprintf("sold %d shares in %s", quantity*-1, game_stock.Stock.Name)
		feed_item.Title = player.User.Name
		feed_item.ImageRoot = player.User.ProfileRoot
	}

	if feedItemMessage == StartGame {
		feed_item.Message = "good luck!"
		feed_item.Title = "game started"
		feed_item.ImageRoot = "/static/imgs/icons/Handshake.svg"
		feed_item.Colour = "grey"
	}

	if feedItemMessage == PlayerPass {
		feed_item.Message = "passed their go"
		feed_item.Title = player.User.Name
		feed_item.ImageRoot = player.User.ProfileRoot
	}

	fmt.Println("creating new FeedItem", feed_item.Message)
	err := db.Create(&feed_item).Error

	if err != nil {
		return FeedItem{}, err
	}

	return feed_item, nil
}
