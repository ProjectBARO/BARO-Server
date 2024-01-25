package services_test

import (
	"errors"
	"gdsc/baro/app/user/models"
	"gdsc/baro/app/user/services"
	"gdsc/baro/app/user/types"
	"gdsc/baro/global/auth"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindOrCreateByEmail(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) (models.User, error) {
	args := m.Called(user)
	return *args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (models.User, error) {
	args := m.Called(id)
	return *args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type MockUserUtil struct {
	mock.Mock
}

func (m *MockUserUtil) FindCurrentUser(c *gin.Context) (*models.User, error) {
	args := m.Called(c)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_Login(t *testing.T) {
	// Mock UserRepository
	mockRepo := new(MockUserRepository)

	// Set up sample user for the test
	input := types.RequestCreateUser{
		Name:     "test",
		Email:    "test@gmail.com",
		Age:      25,
		Gender:   "male",
		FcmToken: "test_token",
	}

	expectedUser := &models.User{
		ID:       1,
		Name:     input.Name,
		Nickname: input.Name,
		Email:    input.Email,
		Age:      input.Age,
		Gender:   input.Gender,
		FcmToken: input.FcmToken,
	}

	// Set up expectations for the mock repository
	mockRepo.On("FindOrCreateByEmail", mock.AnythingOfType("*models.User")).Return(expectedUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(*expectedUser, nil)

	// Create UserService with the mock repository
	userService := services.NewUserService(mockRepo, nil)

	// Call the method under test
	responseToken, err := userService.Login(input)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)

	// Check the results
	assert.Nil(t, err)
	assert.NotEmpty(t, responseToken.Token)
}

func TestUserService_GetUserInfo(t *testing.T) {
	// Mock UserRepository
	mockRepo := new(MockUserRepository)
	mockUtil := new(MockUserUtil)

	// Set up sample user for the test
	expectedUser := &models.User{
		ID:     1,
		Name:   "test",
		Email:  "test@gmail.com",
		Age:    25,
		Gender: "male",
	}

	// Set up expectations for the mock util
	mockUtil.On("FindCurrentUser", mock.AnythingOfType("*gin.Context")).Return(expectedUser, nil)

	// Create UserService with the mock repository and util
	userService := services.NewUserService(mockRepo, mockUtil)

	// Create a test context
	ctx, _ := gin.CreateTestContext(nil)
	ctx.Set(string(auth.UserIDKey), "1")

	// Call the method under test
	responseUser, err := userService.GetUserInfo(ctx)

	// Assert that the expectations were met
	mockUtil.AssertExpectations(t)

	// Check the results
	assert.Nil(t, err)
	assert.Equal(t, expectedUser.ID, responseUser.ID)
	assert.Equal(t, expectedUser.Name, responseUser.Name)
	assert.Equal(t, expectedUser.Email, responseUser.Email)
	assert.Equal(t, expectedUser.Age, responseUser.Age)
}

func TestUserService_GetUserInfo_Error(t *testing.T) {
	// Mock UserRepository
	mockRepo := new(MockUserRepository)
	mockUtil := new(MockUserUtil)

	// Set up expectations for the mock util
	mockUtil.On("FindCurrentUser", mock.AnythingOfType("*gin.Context")).Return((*models.User)(nil), errors.New("user not found"))

	// Create UserService with the mock repository and util
	userService := services.NewUserService(mockRepo, mockUtil)

	// Create a test context
	ctx, _ := gin.CreateTestContext(nil)

	// Call the method under test
	responseUser, err := userService.GetUserInfo(ctx)

	// Assert that the expectations were met
	mockUtil.AssertExpectations(t)

	// Check the results
	assert.NotNil(t, err)
	assert.EqualError(t, err, "user not found")
	assert.Equal(t, types.ResponseUser{}, responseUser)
}

func TestUserService_UpdateUserInfo(t *testing.T) {
	// Mock UserRepository
	mockRepo := new(MockUserRepository)
	mockUtil := new(MockUserUtil)

	// Set up sample user for the test
	input := types.RequestUpdateUser{
		Nickname: "new-nickname",
		Age:      30,
		Gender:   "female",
	}

	currentUser := &models.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      25,
		Gender:   "male",
	}

	updatedUser := &models.User{
		ID:       1,
		Name:     "test",
		Nickname: "new-nickname",
		Email:    "test@gmail.com",
		Age:      30,
		Gender:   "female",
	}

	// Set up expectations for the mock repository and util
	mockUtil.On("FindCurrentUser", mock.AnythingOfType("*gin.Context")).Return(currentUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(*updatedUser, nil)

	// Create UserService with the mock repository and util
	userService := services.NewUserService(mockRepo, mockUtil)

	// Create a test context
	ctx, _ := gin.CreateTestContext(nil)

	// Call the method under test
	responseUser, err := userService.UpdateUserInfo(ctx, input)

	// Assert that the expectations were met
	mockUtil.AssertExpectations(t)
	mockRepo.AssertExpectations(t)

	// Check the results
	assert.Nil(t, err)
	assert.Equal(t, currentUser.ID, responseUser.ID)
	assert.Equal(t, currentUser.Name, responseUser.Name)
	assert.Equal(t, input.Nickname, responseUser.Nickname)
	assert.Equal(t, currentUser.Email, responseUser.Email)
	assert.Equal(t, input.Age, responseUser.Age)
	assert.Equal(t, input.Gender, responseUser.Gender)
}

func TestUserService_UpdateUserInfo_Error(t *testing.T) {
	// Mock UserRepository
	mockRepo := new(MockUserRepository)
	mockUtil := new(MockUserUtil)

	// Set up sample user for the test
	input := types.RequestUpdateUser{
		Nickname: "new-nickname",
		Age:      30,
		Gender:   "female",
	}

	currentUser := &models.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      25,
		Gender:   "male",
	}

	// Set up expectations for the mock util
	mockUtil.On("FindCurrentUser", mock.AnythingOfType("*gin.Context")).Return(currentUser, errors.New("user not found"))

	// Create UserService with the mock repository and util
	userService := services.NewUserService(mockRepo, mockUtil)

	// Create a test context
	ctx, _ := gin.CreateTestContext(nil)

	// Call the method under test
	_, err := userService.UpdateUserInfo(ctx, input)

	// Assert that the expectations were met
	mockUtil.AssertExpectations(t)

	// Check the results
	assert.NotNil(t, err)
	assert.EqualError(t, err, "user not found")
}

func TestUserService_DeleteUser(t *testing.T) {
	// Mock UserRepository
	mockRepo := new(MockUserRepository)
	mockUtil := new(MockUserUtil)

	// Set up sample user for the test
	currentUser := &models.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      25,
		Gender:   "male",
	}

	// Set up expectations for the mock repository and util
	mockUtil.On("FindCurrentUser", mock.AnythingOfType("*gin.Context")).Return(currentUser, nil)
	mockRepo.On("Delete", currentUser).Return(nil)

	// Create UserService with the mock repository and util
	userService := services.NewUserService(mockRepo, mockUtil)

	// Create a test context
	ctx, _ := gin.CreateTestContext(nil)

	// Call the method under test
	err := userService.DeleteUser(ctx)

	// Assert that the expectations were met
	mockUtil.AssertExpectations(t)
	mockRepo.AssertExpectations(t)

	// Check the results
	assert.Nil(t, err)
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	// Mock UserRepository
	mockRepo := new(MockUserRepository)
	mockUtil := new(MockUserUtil)

	// Set up sample user for the test
	currentUser := &models.User{
		ID:       1,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      25,
		Gender:   "male",
	}

	// Set up expectations for the mock util
	mockUtil.On("FindCurrentUser", mock.AnythingOfType("*gin.Context")).Return(currentUser, errors.New("user not found"))

	// Create UserService with the mock repository and util
	userService := services.NewUserService(mockRepo, mockUtil)

	// Create a test context
	ctx, _ := gin.CreateTestContext(nil)

	// Call the method under test
	err := userService.DeleteUser(ctx)

	// Assert that the expectations were met
	mockUtil.AssertExpectations(t)

	// Check the results
	assert.NotNil(t, err)
	assert.EqualError(t, err, "user not found")
}
