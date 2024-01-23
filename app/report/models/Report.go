package models

import "time"

type Report struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint
	AlertCount   int
	AnalysisTime int
	Type         string
	Predict      string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
