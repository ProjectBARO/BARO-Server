package services

import (
	"fmt"
	"gdsc/baro/app/user/models"
	"gdsc/baro/app/user/repositories"
	"gdsc/baro/app/user/types"
	"gdsc/baro/global/auth"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	UserRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (service *UserService) generateToken(userID uint) (string, error) {
	tokenClaim := auth.NewClaim(fmt.Sprint(userID))
	return auth.GenerateToken(tokenClaim)
}

func (service *UserService) Login(input types.RequestCreateUser) (types.ResponseToken, error) {
	requestCreateUser := models.User{
		Name:     input.Name,
		Nickname: input.Name,
		Email:    input.Email,
		Age:      input.Age,
		Gender:   input.Gender,
	}

	user, err := service.UserRepository.FindOrCreateByEmail(&requestCreateUser)
	if err != nil {
		return types.ResponseToken{}, err
	}

	token, err := service.generateToken(user.ID)
	if err != nil {
		return types.ResponseToken{}, err
	}

	return types.ResponseToken{Token: token}, nil
}

func (service *UserService) GetUserInfo(c *gin.Context) (types.ResponseUser, error) {
	user, err := service.FindCurrentUser(c)
	if err != nil {
		return types.ResponseUser{}, err
	}

	responseUser := types.ResponseUser{
		ID:       user.ID,
		Name:     user.Name,
		Nickname: user.Nickname,
		Email:    user.Email,
		Age:      user.Age,
		Gender:   user.Gender,
	}
	return responseUser, nil
}

func (service *UserService) UpdateUserInfo(c *gin.Context, input types.RequestUpdateUser) (types.ResponseUser, error) {
	user, err := service.FindCurrentUser(c)
	if err != nil {
		return types.ResponseUser{}, err
	}

	service.updateUser(user, input)

	updatedUser, updateErr := service.UserRepository.Update(user)
	if updateErr != nil {
		return types.ResponseUser{}, updateErr
	}

	responseUser := types.ResponseUser{
		ID:       updatedUser.ID,
		Name:     updatedUser.Name,
		Nickname: updatedUser.Nickname,
		Email:    updatedUser.Email,
		Age:      updatedUser.Age,
		Gender:   updatedUser.Gender,
	}

	return responseUser, nil
}

func (service *UserService) updateUser(user *models.User, input types.RequestUpdateUser) {
	if input.Nickname != "" {
		user.Nickname = input.Nickname
	}

	if input.Age != 0 {
		user.Age = input.Age
	}

	if input.Gender != "" {
		user.Gender = input.Gender
	}
}

func (service *UserService) DeleteUser(c *gin.Context) error {
	user, err := service.FindCurrentUser(c)
	if err != nil {
		return err
	}

	return service.UserRepository.Delete(user)
}

func (service *UserService) FindCurrentUser(c *gin.Context) (*models.User, error) {
	userID, exists := c.Get(string(auth.UserIDKey))
	if !exists {
		return nil, fmt.Errorf("not found user id")
	}

	user, err := service.UserRepository.FindByID(userID.(string))
	if err != nil {
		return nil, fmt.Errorf("not found user")
	}

	return &user, nil
}
