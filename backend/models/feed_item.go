package models

import (
	"fmt"

	"gorm.io/gorm"
)

type FeedItem struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`

	Message string

	GameStock   GameStock
	GameStockID uint
	Player      Player
	PlayerID    uint

	Game   Game
	GameID string
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

func NewFeedItem(playerStock PlayerStock, game Game, quantity int, db *gorm.DB) (FeedItem, error) {

	feed_item := FeedItem{
		GameStock: playerStock.GameStock,
		Player:    playerStock.Player,
		Game:      game,
	}

	if quantity > 0 {
		feed_item.Message = fmt.Sprintf("bought %d shares in %s", quantity, playerStock.GameStock.Stock.Name)
	} else {
		feed_item.Message = fmt.Sprintf("sold %d shares in %s", quantity*-1, playerStock.GameStock.Stock.Name)
	}

	fmt.Println("creating new FeedItem", feed_item)
	err := db.Create(&feed_item).Error

	if err != nil {
		return FeedItem{}, err
	}

	return feed_item, nil
}
