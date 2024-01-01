package services

import (
	"gdsc/baro/models/repositories"
	"gdsc/baro/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoService struct {
	VideoRepository *repositories.VideoRepository
}

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{
		VideoRepository: repositories.NewVideoRepository(db),
	}
}

func (service *VideoService) GetVideos(c *gin.Context) ([]types.ResponseVideo, error) {
	videos, err := service.VideoRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var responseVideos []types.ResponseVideo
	for _, video := range videos {
		responseVideos = append(responseVideos, types.ResponseVideo{
			VideoID:      video.VideoID,
			Title:        video.Title,
			ThumbnailUrl: video.ThumbnailUrl,
		})
	}

	return responseVideos, nil
}
