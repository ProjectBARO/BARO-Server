package services_test

import (
	"errors"
	"gdsc/baro/app/video/models"
	"gdsc/baro/app/video/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVideoRepository struct {
	mock.Mock
}

func (m *MockVideoRepository) FindAll() ([]models.Video, error) {
	args := m.Called()
	return args.Get(0).([]models.Video), args.Error(1)
}

func (m *MockVideoRepository) FindByCategory(category string) ([]models.Video, error) {
	args := m.Called(category)
	return args.Get(0).([]models.Video), args.Error(1)
}

func TestVideoService_GetVideos(t *testing.T) {
	// Mock VideoRepository
	mockRepo := new(MockVideoRepository)

	// Set up sample video for the test
	expectedVideos := []models.Video{
		{VideoID: "1", Title: "Video 1", ThumbnailUrl: "thumbnail1.jpg", Category: "category1"},
		{VideoID: "2", Title: "Video 2", ThumbnailUrl: "thumbnail2.jpg", Category: "category2"},
	}

	// Set up expectations for the mock repository
	mockRepo.On("FindAll").Return(expectedVideos, nil)

	// Create VideoService with the mock repository
	videoService := services.NewVideoService(mockRepo)

	// Create a test context
	ctx, _ := gin.CreateTestContext(nil)

	// Call the method under test
	responseVideos, err := videoService.GetVideos(ctx)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)

	// Check the results
	assert.Nil(t, err)
	assert.Equal(t, len(expectedVideos), len(responseVideos))
	assert.Equal(t, expectedVideos[0].Title, responseVideos[0].Title)
	assert.Equal(t, expectedVideos[1].ThumbnailUrl, responseVideos[1].ThumbnailUrl)
}

func TestVideoService_GetVideos_ErrorInRepository(t *testing.T) {
	// Mock VideoRepository
	mockRepo := new(MockVideoRepository)

	// Set up expectations for the mock repository to return an error
	mockRepo.On("FindAll").Return([]models.Video{}, errors.New("database error"))

	// Create VideoService with the mock repository
	videoService := services.NewVideoService(mockRepo)

	// Create a test context
	ctx, _ := gin.CreateTestContext(nil)

	// Call the method under test
	responseVideos, err := videoService.GetVideos(ctx)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)

	// Check the results
	assert.Error(t, err)
	assert.Nil(t, responseVideos)
}

func TestVideoService_GetVideosByCategory(t *testing.T) {
	// Mock VideoRepository
	mockRepo := new(MockVideoRepository)

	// Set up sample video for the test
	expectedVideos := []models.Video{
		{VideoID: "1", Title: "Video 1", ThumbnailUrl: "thumbnail1.jpg", Category: "category1"},
	}

	// Set up expectations for the mock repository
	mockRepo.On("FindByCategory", mock.Anything).Return(expectedVideos, nil)

	// Create VideoService with the mock repository
	videoService := services.NewVideoService(mockRepo)

	// Create a test context
	ctx, _ := gin.CreateTestContext(nil)

	// Call the method under test
	responseVideos, err := videoService.GetVideosByCategory(ctx)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)

	// Check the results
	assert.Nil(t, err)
	assert.Equal(t, len(expectedVideos), len(responseVideos))
	assert.Equal(t, expectedVideos[0].Title, responseVideos[0].Title)
	assert.Equal(t, expectedVideos[0].ThumbnailUrl, responseVideos[0].ThumbnailUrl)
}
