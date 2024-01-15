package models

// Define a GORM model
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Phone string `gorm:"not null"`
	Email string `gorm:"not null"`
}
