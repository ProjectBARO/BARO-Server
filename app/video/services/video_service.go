package services

import (
	"gdsc/baro/app/video/repositories"
	"gdsc/baro/app/video/types"

	"github.com/gin-gonic/gin"
)

type VideoService struct {
	VideoRepository *repositories.VideoRepository
}

func NewVideoService(videoRepository *repositories.VideoRepository) *VideoService {
	return &VideoService{
		VideoRepository: videoRepository,
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
