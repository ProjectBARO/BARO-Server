package utils_test

import (
	"fmt"
	"gdsc/baro/app/user/models"
	"gdsc/baro/global/auth"
	"gdsc/baro/global/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepositoryInterface
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

func TestUserUtil_FindCurrentUser(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := new(MockUserRepository)

	// Create an instance of UserUtil with the mock repository
	userUtil := utils.NewUserUtil(mockRepo)

	// Create a sample user ID for testing
	userID := "123"

	// Create a sample gin.Context with the user ID
	ctx := &gin.Context{}
	ctx.Set(string(auth.UserIDKey), userID)

	// Create a sample user for testing
	expectedUser := &models.User{
		ID:       123,
		Name:     "test",
		Nickname: "test",
		Email:    "test@gmail.com",
		Age:      25,
		Gender:   "male",
	}

	// Set up expectations for the mock repository
	mockRepo.On("FindByID", userID).Return(expectedUser, nil)

	// Call the method under test
	user, err := userUtil.FindCurrentUser(ctx)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)

	// Validate the result
	assert.Nil(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserUtil_FindCurrentUser_UserNotFound(t *testing.T) {
	// Create a new instance of the mock repository
	mockRepo := new(MockUserRepository)

	// Create an instance of UserUtil with the mock repository
	userUtil := utils.NewUserUtil(mockRepo)

	// Create a sample user ID for testing
	userID := "123"

	// Create a sample gin.Context with the user ID
	ctx := &gin.Context{}
	ctx.Set(string(auth.UserIDKey), userID)

	// Set up expectations for the mock repository when user is not found
	mockRepo.On("FindByID", userID).Return(&models.User{}, fmt.Errorf("user not found"))

	// Call the method under test
	user, err := userUtil.FindCurrentUser(ctx)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)

	// Validate the result
	assert.NotNil(t, err)
	assert.Nil(t, user)
}
