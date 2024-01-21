package services

import (
	"gdsc/baro/app/video/repositories"
	"gdsc/baro/app/video/types"

	"github.com/gin-gonic/gin"
)

type VideoServiceInterface interface {
	GetVideos(c *gin.Context) ([]types.ResponseVideo, error)
	GetVideosByCategory(c *gin.Context) ([]types.ResponseVideo, error)
}

type VideoService struct {
	VideoRepository repositories.VideoRepositoryInterface
}

func NewVideoService(videoRepository repositories.VideoRepositoryInterface) *VideoService {
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
			Category:     video.Category,
		})
	}

	return responseVideos, nil
}

func (service *VideoService) GetVideosByCategory(c *gin.Context) ([]types.ResponseVideo, error) {
	category := c.Query("keyword")

	videos, err := service.VideoRepository.FindByCategory(category)
	if err != nil {
		return nil, err
	}

	var responseVideos []types.ResponseVideo
	for _, video := range videos {
		responseVideos = append(responseVideos, types.ResponseVideo{
			VideoID:      video.VideoID,
			Title:        video.Title,
			ThumbnailUrl: video.ThumbnailUrl,
			Category:     video.Category,
		})
	}

	return responseVideos, nil
}
