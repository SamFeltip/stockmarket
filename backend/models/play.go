package models

type Play struct {
	ID          uint `gorm:"primaryKey"`
	GameStock   GameStock
	GameStockID uint
	Player      Player
	PlayerID    uint

	Game   Game
	GameID string
}

NewPlay(playerStock PlayerStock, quantity int, db *gorm.DB) (Play, error) {
	play := Play{
		GameStock:   playerStock.GameStock,
		Player: 	playerStock.Player,
	}
