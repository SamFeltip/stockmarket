package models

import "gorm.io/gorm"

type GameDisplay struct {
	ID              string
	PlayerCount     int
	PeriodCount     int
	CurrentPeriod   int
	CurrentUserName string
	CurrentUserID   uint
	FeedItems       []FeedItem
	Status          string
}

func LoadGameDisplay(gameID string, db *gorm.DB) (GameDisplay, error) {

	var gameDisplay GameDisplay

	err := db.Table("games as g").
		Select("g.id, count(p.id) as player_count, g.period_count, g.current_period, u.name as current_user_name, g.current_user_id, fi.id as feed_item_id, g.status as status").
		Joins("inner join players as p on p.game_id = g.id").
		Joins("inner join users as u on g.current_user_id = u.id").
		Joins("inner join feed_items as fi on g.id = fi.game_id").
		Where("g.id = ? and fi.colour = ? and g.current_period = fi.period", "dogs", "").
		Group("g.current_user_id, g.id, u.name, fi.id, fi.colour").
		Find(&gameDisplay).Error
	return gameDisplay, err
}

func (gameDisplay *GameDisplay) CurrentTurn() int {

	var currentPeriodPlays []FeedItem
	for _, feedItem := range gameDisplay.FeedItems {
		if feedItem.Period == gameDisplay.CurrentPeriod {
			currentPeriodPlays = append(currentPeriodPlays, feedItem)
		}
	}

	return len(currentPeriodPlays)
}
