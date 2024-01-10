package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gdsc/baro/app/video/controllers"
	"gdsc/baro/app/video/types"
)

// MockVideoService is a mock implementation of VideoServiceInterface
type MockVideoService struct {
	mock.Mock
}

func (m *MockVideoService) GetVideos(c *gin.Context) ([]types.ResponseVideo, error) {
	args := m.Called(c)
	return args.Get(0).([]types.ResponseVideo), args.Error(1)
}

// Because the data of global.Response is of the interface{} type
// I created a Response structure for testing.
type Response struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    []types.ResponseVideo `json:"data"`
}

func TestVideoController_GetVideos(t *testing.T) {
	// Mock VideoService
	mockService := new(MockVideoService)

	// Create a sample video for the test
	expectedVideos := []types.ResponseVideo{
		{VideoID: "1", Title: "Video 1", ThumbnailUrl: "thumbnail1.jpg"},
		{VideoID: "2", Title: "Video 2", ThumbnailUrl: "thumbnail2.jpg"},
	}

	// Set up expectations for the mock service
	mockService.On("GetVideos", mock.Anything).Return(expectedVideos, nil)

	// Create VideoController with the mock service
	videoController := controllers.NewVideoController(mockService)

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/videos", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Call the method under test
	videoController.GetVideos(ctx)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)

	// Create the expected response
	expectedJson := Response{
		Status:  200,
		Message: "success",
		Data:    expectedVideos,
	}

	// Create the actual response
	var actualJson Response
	err := json.Unmarshal(w.Body.Bytes(), &actualJson)
	if err != nil {
		fmt.Println(err)
	}

	// Check the response
	assert.Equal(t, 200, w.Code)

	// Check the response body
	assert.Equal(t, expectedJson, actualJson)
	assert.Equal(t, expectedJson.Status, actualJson.Status)
	assert.Equal(t, expectedJson.Message, actualJson.Message)
	assert.Equal(t, expectedJson.Data, actualJson.Data)
}

func TestVideoController_GetVideos_ErrorInService(t *testing.T) {
	// Mock VideoService
	mockService := new(MockVideoService)

	// Set up expectations for the mock service to return an error
	mockService.On("GetVideos", mock.Anything).Return([]types.ResponseVideo{}, errors.New("failed"))

	// Create VideoController with the mock service
	videoController := controllers.NewVideoController(mockService)

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/videos", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Call the method under test
	videoController.GetVideos(ctx)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)

	// Create the expected response
	expectedJson := Response{
		Status:  400,
		Message: "failed",
		Data:    nil,
	}

	// Create the actual response
	var actualJson Response
	err := json.Unmarshal(w.Body.Bytes(), &actualJson)
	if err != nil {
		fmt.Println(err)
	}

	// Check the response
	assert.Equal(t, 400, w.Code)

	// Check the response body
	assert.Equal(t, expectedJson, actualJson)
	assert.Equal(t, expectedJson.Status, actualJson.Status)
	assert.Equal(t, expectedJson.Message, actualJson.Message)
	assert.Equal(t, expectedJson.Data, actualJson.Data)
}
