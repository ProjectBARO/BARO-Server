package models

import "time"

type Report struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Predict   string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
