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
	existingUser, err := service.UserRepository.FindByEmail(input.Email)
	if err != nil {
		return types.ResponseToken{}, err
	}

	if existingUser != nil {
		token, err := service.generateToken(existingUser.ID)
		if err != nil {
			return types.ResponseToken{}, err
		}

		return types.ResponseToken{Token: token}, nil
	}

	token, err := service.RegisterUser(input)
	if err != nil {
		return types.ResponseToken{}, err
	}

	return token, nil
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

func (service *UserService) FindUserByID(c *gin.Context) (types.ResponseUser, error) {
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
