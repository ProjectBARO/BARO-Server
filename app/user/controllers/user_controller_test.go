package controllers_test

import (
	"encoding/json"
	"errors"
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

func (m *MockUserService) UpdateFcmToken(c *gin.Context, input types.RequestUpdateFcmToken) error {
	args := m.Called(c, input)
	return args.Error(0)
}

func (m *MockUserService) GetUserInfo(c *gin.Context) (types.ResponseUser, error) {
	args := m.Called(c)
	if args.Get(0) == nil {
		return types.ResponseUser{}, args.Error(1)
	}
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
		Name:     "test",
		Email:    "test@gmail.com",
		Age:      20,
		Gender:   "male",
		FcmToken: "test_token",
	}

	// Set up expectations for the mock service
	mockService.On("Login", input).Return(types.ResponseToken{}, nil)

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(mockRequestBody(input))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.LoginOrRegisterUser(c)

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

func TestUserController_LoginOrRegisterUser_InvalidJson(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Create a test context using httptest
	w := httptest.NewRecorder()

	// Create a test context using httptest
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(strings.NewReader("invalid json"))

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.LoginOrRegisterUser(c)

	// Validate the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Parse the response body
	var response ResponseToken
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Check the results
	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Equal(t, "invalid character 'i' looking for beginning of value", response.Message)
	assert.Equal(t, types.ResponseToken{}, response.Data)
}

func TestUserController_LoginOrRegisterUser_InvalidInput_Email(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Create a sample request for the test
	input := types.RequestCreateUser{
		Name:     "test",
		Email:    "test",
		Age:      20,
		Gender:   "male",
		FcmToken: "test_token",
	}

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", nil)
	req.Body = io.NopCloser(mockRequestBody(input))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.LoginOrRegisterUser(c)

	// Validate the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Parse the response body
	var response ResponseToken
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Check the results
	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Equal(t, "Key: 'RequestCreateUser.Email' Error:Field validation for 'Email' failed on the 'email' tag", response.Message)
	assert.Equal(t, types.ResponseToken{}, response.Data)
}

func TestUserController_FcmTokenUpdate(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Create a sample request for the test
	input := types.RequestUpdateFcmToken{
		FcmToken: "new-fcm-token",
	}

	// Set up expectations for the mock service
	mockService.On("UpdateFcmToken", mock.Anything, input).Return(nil)

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/me/fcm-token", nil)
	req.Body = io.NopCloser(mockRequestBody(input))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.UpdateFcmToken(c)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)

	// Validate the response
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_FcmTokenUpdate_InvalidJson(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Create a test context using httptest
	w := httptest.NewRecorder()

	// Create a test context using httptest
	req, _ := http.NewRequest("PUT", "/users/me/fcm-token", nil)
	req.Body = io.NopCloser(strings.NewReader("invalid json"))

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.UpdateFcmToken(c)

	// Validate the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Parse the response body
	var response ResponseToken
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Check the results
	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Equal(t, "invalid character 'i' looking for beginning of value", response.Message)
	assert.Equal(t, types.ResponseToken{}, response.Data)
}

func TestUserController_FcmTokenUpdate_InvalidInput_FcmToken(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Create a sample request for the test
	input := types.RequestUpdateFcmToken{
		FcmToken: "",
	}

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/me/fcm-token", nil)
	req.Body = io.NopCloser(mockRequestBody(input))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.UpdateFcmToken(c)

	// Validate the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Parse the response body
	var response ResponseToken
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Check the results
	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Equal(t, "Key: 'RequestUpdateFcmToken.FcmToken' Error:Field validation for 'FcmToken' failed on the 'required' tag", response.Message)
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.GetUserInfo(c)

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

func TestUserController_GetUserInfo_Error(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)

	// Set up expectations for the mock service
	mockService.On("GetUserInfo", mock.Anything).Return(types.ResponseUser{}, errors.New("user not found"))

	// Create a test context using httptest
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/me", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.GetUserInfo(c)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)

	// Validate the response
	assert.Equal(t, 400, w.Code)

	// Parse the response body
	var response ResponseUser
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Check the results
	assert.Equal(t, 400, response.Status)
	assert.Equal(t, "user not found", response.Message)
	assert.Equal(t, types.ResponseUser{}, response.Data)
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.UpdateUserInfo(c)

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

func TestUserController_UpdateUserInfo_InvalidJson(t *testing.T) {
	// Mock UserService
	mockService := new(MockUserService)

	// Create UserController with the mock service
	userController := controllers.NewUserController(mockService)
	// Create a test context using httptest
	w := httptest.NewRecorder()

	// Create a test context using httptest
	req, _ := http.NewRequest("PUT", "/users/me", nil)
	req.Body = io.NopCloser(strings.NewReader("invalid json"))

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.UpdateUserInfo(c)

	// Validate the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Parse the response body
	var response ResponseUser
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Check the results
	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Equal(t, "invalid character 'i' looking for beginning of value", response.Message)
	assert.Equal(t, types.ResponseUser{}, response.Data)
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
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the method under test
	userController.DeleteUser(c)

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
