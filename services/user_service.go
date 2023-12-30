package services

import (
	"fmt"
	"gdsc/baro/auth"
	"gdsc/baro/models"
	"gdsc/baro/models/repositories"
	"gdsc/baro/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository *repositories.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		UserRepository: repositories.NewUserRepository(db),
	}
}

func (service *UserService) generateToken(userID uint) (string, error) {
	tokenClaim := auth.NewClaim(fmt.Sprint(userID))
	return auth.GenerateToken(tokenClaim)
}

func (service *UserService) Login(input types.RequestCreateUser) (types.ResponseToken, error) {
	existingUser, dbErr := service.UserRepository.FindByEmail(input.Email)
	if dbErr != nil && dbErr.Error() == "record not found" {
		responseToken, registerErr := service.RegisterUser(input)
		if registerErr != nil {
			return types.ResponseToken{}, registerErr
		}

		return responseToken, nil
	}

	if existingUser != nil {
		token, err := service.generateToken(existingUser.ID)
		if err != nil {
			return types.ResponseToken{}, err
		}

		responseToken := types.ResponseToken{Token: token}
		return responseToken, nil
	}

	return types.ResponseToken{}, dbErr
}

func (service *UserService) RegisterUser(input types.RequestCreateUser) (types.ResponseToken, error) {
	userToCreate := models.User{
		Name:     input.Name,
		Nickname: input.Name,
		Email:    input.Email,
		Age:      input.Age,
		Gender:   input.Gender,
	}

	newUser, err := service.UserRepository.Create(&userToCreate)
	if err != nil {
		return types.ResponseToken{}, err
	}

	token, err := service.generateToken(newUser.ID)
	if err != nil {
		return types.ResponseToken{}, err
	}

	return types.ResponseToken{Token: token}, nil
}

func (service *UserService) GetUserInfo(c *gin.Context) (types.ResponseUser, error) {
	user := auth.FindCurrentUser(c)

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
	user := auth.FindCurrentUser(c)

	service.updateUser(user, input)

	updatedUser, err := service.UserRepository.Update(user)
	if err != nil {
		return types.ResponseUser{}, err
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
	user := auth.FindCurrentUser(c)

	return service.UserRepository.Delete(user)
}
