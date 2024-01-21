package models

type Video struct {
	VideoID      string `gorm:"primaryKey"`
	Title        string
	ThumbnailUrl string
	Category     string
}
