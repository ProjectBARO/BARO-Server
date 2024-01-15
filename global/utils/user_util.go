package utils

import (
	"fmt"
	"gdsc/baro/app/user/models"
	"gdsc/baro/app/user/repositories"
	"gdsc/baro/global/auth"

	"github.com/gin-gonic/gin"
)

type UserUtilInterface interface {
	FindCurrentUser(c *gin.Context) (*models.User, error)
}

type UserUtil struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUserUtil(userRepository repositories.UserRepositoryInterface) *UserUtil {
	return &UserUtil{
		UserRepository: userRepository,
	}
}

func (util *UserUtil) FindCurrentUser(c *gin.Context) (*models.User, error) {
	userID, exists := c.Get(string(auth.UserIDKey))
	if !exists {
		return nil, fmt.Errorf("not found user id")
	}

	user, err := util.UserRepository.FindByID(userID.(string))
	if err != nil {
		return nil, fmt.Errorf("not found user")
	}

	return &user, nil
}
