package pb

import (
	"context"
	"gdsc/baro/app/video/repositories"

	videopb "gdsc/baro/protos/video"
)

type VideoPbApp struct {
	VideoRepository repositories.VideoRepositoryInterface
	videopb.UnimplementedVideoServiceServer
}

func NewVideoPbApp(videoRepository repositories.VideoRepositoryInterface) *VideoPbApp {
	return &VideoPbApp{
		VideoRepository: videoRepository,
	}
}

func (s *VideoPbApp) GetVideos(ctx context.Context, req *videopb.GetVideosRequest) (*videopb.VideosResponse, error) {
	videos, err := s.VideoRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var responseVideos []*videopb.VideoMessage
	for _, video := range videos {
		responseVideos = append(responseVideos, &videopb.VideoMessage{
			VideoId:      video.VideoID,
			Title:        video.Title,
			ThumbnailUrl: video.ThumbnailUrl,
			Category:     video.Category,
		})
	}

	return &videopb.VideosResponse{
		Videos: responseVideos,
	}, nil
}

func (s *VideoPbApp) GetVideosByCategory(ctx context.Context, req *videopb.GetVideosByCategoryRequest) (*videopb.VideosResponse, error) {
	category := req.GetCategory()

	videos, err := s.VideoRepository.FindByCategory(category)
	if err != nil {
		return nil, err
	}

	var responseVideos []*videopb.VideoMessage
	for _, video := range videos {
		responseVideos = append(responseVideos, &videopb.VideoMessage{
			VideoId:      video.VideoID,
			Title:        video.Title,
			ThumbnailUrl: video.ThumbnailUrl,
			Category:     video.Category,
		})
	}

	return &videopb.VideosResponse{
		Videos: responseVideos,
	}, nil
}

func (s *VideoPbApp) mustEmbedUnimplementedVideoServiceServer() {}
