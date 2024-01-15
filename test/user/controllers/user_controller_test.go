package controllers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gdsc/baro/app/user/controllers"
	"gdsc/baro/app/user/types"
)

// MockUserService is a mock implementation of UserServiceInterface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(input types.RequestCreateUser) (types.ResponseToken, error) {
	args := m.Called(input)
	return args.Get(0).(types.ResponseToken), args.Error(1)
}

func (m *MockUserService) GetUserInfo(c *gin.Context) (types.ResponseUser, error) {
	args := m.Called(c)
	return args.Get(0).(types.ResponseUser), args.Error(1)
}

func (m *MockUserService) UpdateUserInfo(c *gin.Context, input types.RequestUpdateUser) (types.ResponseUser, error) {
	args := m.Called(c, input)
	return args.Get(0).(types.ResponseUser), args.Error(1)
}

func (m *MockUserService) DeleteUser(c *gin.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

// Because the data of global.Response is of the interface{} type
// I created a Response structure for testing.
type ResponseToken struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    types.ResponseToken `json:"data"`
}

type ResponseUser struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    types.ResponseUser `json:"data"`
}

func TestUserController_LoginOrRegisterUser(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Create a sample request for the test
	input := types.RequestCreateUser{
		Name:   "test",
		Email:  "test@gmail.com",
		Age:    20,
		Gender: "male",
	}

	// Set up expectations for the mock service
	mockService.On("Login", input).Return(types.ResponseToken{}, nil)

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(mockRequestBody(input))
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Call the method under test
	userController.LoginOrRegisterUser(ctx)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)

	// Validate the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body
	var response ResponseToken
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Check the results
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "success", response.Message)
	assert.Equal(t, types.ResponseToken{}, response.Data)
}

func TestUserController_GetUserInfo(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Create a sample user for the test
	expectedUser := types.ResponseUser{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
	}

	// Set up expectations for the mock service
	mockService.On("GetUserInfo", mock.Anything).Return(expectedUser, nil)

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/me", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Call the method under test
	userController.GetUserInfo(ctx)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)

	// Validate the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body
	var response ResponseUser
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Check the results
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "success", response.Message)
	assert.Equal(t, expectedUser, response.Data)
}

func TestUserController_UpdateUserInfo(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Create a sample request for the test
	input := types.RequestUpdateUser{
		Nickname: "new-nickname",
		Age:      30,
		Gender:   "female",
	}

	updatedUser := types.ResponseUser{
		ID:       1,
		Name:     "test",
		Nickname: "new-nickname",
		Email:    "test@gmail.com",
		Age:      30,
		Gender:   "female",
	}

	// Set up expectations for the mock service
	mockService.On("UpdateUserInfo", mock.Anything, input).Return(updatedUser, nil)

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/me", mockRequestBody(input))
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Call the method under test
	userController.UpdateUserInfo(ctx)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)

	// Validate the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body
	var response ResponseUser
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Check the results
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "success", response.Message)
	assert.Equal(t, updatedUser, response.Data)
}

func TestUserController_DeleteUser(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Set up expectations for the mock service
	mockService.On("DeleteUser", mock.Anything).Return(nil)

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/me", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Call the method under test
	userController.DeleteUser(ctx)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)

	// Validate the response
	assert.Equal(t, http.StatusOK, w.Code)
}

// Helper function to convert an object to JSON
func mockRequestBody(v interface{}) *strings.Reader {
	raw, _ := json.Marshal(v)
	return strings.NewReader(string(raw))
}
