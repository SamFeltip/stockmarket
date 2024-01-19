package models

type Player struct {
	GameID int
	Game   Game `gorm:"foreignkey:GameID"`
	UserID uint
	User   User `gorm:"foreignkey:UserID"`
}
