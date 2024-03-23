package models

import "gorm.io/gorm"

type GameDisplay struct {
	ID              string
	PlayerCount     int
	PeriodCount     int
	CurrentPeriod   int
	CurrentUserName string
	CurrentUserID   uint
	CurrentTurn     int
	Status          string
}

func LoadGameDisplay(gameID string, db *gorm.DB) (GameDisplay, error) {

	var gameDisplay GameDisplay

	err := db.Table("games as g").
		Select("g.id, count(p.id) as player_count, g.period_count, g.current_period, u.name as current_user_name, g.current_user_id, g.status as status").
		Joins("inner join players as p on p.game_id = g.id").
		Joins("inner join users as u on g.current_user_id = u.id").
		Where("g.id = ?", gameID).
		Group("g.current_user_id, g.id, u.name").
		Find(&gameDisplay).Error

	if err != nil {
		return gameDisplay, err
	}

	err = gameDisplay.SetCurrentTurn(db)

	return gameDisplay, err
}

func (gameDisplay *GameDisplay) SetCurrentTurn(db *gorm.DB) error {

	// select count(*) from feed_items where game_id = ? and period = ? order by created_at desc
	// limit 1

	err := db.Table("feed_items").
		Select("count(*)").
		Where("game_id = ? and period = ?", gameDisplay.ID, gameDisplay.CurrentPeriod).
		Scan(&gameDisplay.CurrentTurn).Error

	return err
}
