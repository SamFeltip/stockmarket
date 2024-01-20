package models

type Player struct {
	ID       uint
	GameID   string
	Game     Game `gorm:"foreignkey:GameID"`
	UserID   uint
	User     User `gorm:"foreignkey:UserID"`
	Position int
}
