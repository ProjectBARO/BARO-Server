package controllers

import (
	"gdsc/baro/services"
	"gdsc/baro/types"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (controller *UserController) LoginOrRegisterUser(c *gin.Context) {
	var input types.RequestCreateUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := controller.UserService.Login(input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

func (controller *UserController) FindUserByID(c *gin.Context) {
	user, err := controller.UserService.FindUserByID(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}
