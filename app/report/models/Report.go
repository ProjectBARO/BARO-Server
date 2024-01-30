package models

import "time"

type Report struct {
	ID                uint `gorm:"primaryKey"`
	UserID            uint
	AlertCount        int
	AnalysisTime      int
	Type              string
	Predict           string
	NormalRatio       string
	Score             string
	NeckAngles        string
	Distances         string
	StatusFrequencies string
	CreatedAt         time.Time `gorm:"autoCreateTime"`
}
