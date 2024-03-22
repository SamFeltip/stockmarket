package models

import (
	gorm "gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null;unique"`
	StartingValue float64
	ImagePath     string `gorm:"not null;unique"`
	Insights      []Insight
	Variation     float64 // +/- maximum value of variation (0.50 increments)
}
