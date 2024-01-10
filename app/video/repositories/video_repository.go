package repositories

import (
	"gdsc/baro/app/video/models"

	"gorm.io/gorm"
)

type VideoRepositoryInterface interface {
	FindAll() ([]models.Video, error)
}

type VideoRepository struct {
	DB *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{
		DB: db,
	}
}

func (r *VideoRepository) FindAll() ([]models.Video, error) {
	var videos []models.Video

	result := r.DB.Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}

	return videos, nil
}
