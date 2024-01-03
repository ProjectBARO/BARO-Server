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

// @Tags Users
// @Summary Login or register user
// @Description Log in if the user exists, if not, register a new user
// @Accept  json
// @Produce  json
// @Param   user    body    types.RequestCreateUser   true    "user info to login or register"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Router /login [post]
func (controller *UserController) LoginOrRegisterUser(c *gin.Context) {
	var input types.RequestCreateUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	token, err := controller.UserService.Login(input)
	if err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, types.Response{
		Status:  200,
		Message: "success",
		Data:    token,
	})
}

// @Tags Users
// @Summary Get user information
// @Description Get information about the currently logged in user
// @Accept  json
// @Produce  json
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Security Bearer
// @Router /users/me [get]
func (controller *UserController) GetUserInfo(c *gin.Context) {
	user, err := controller.UserService.GetUserInfo(c)
	if err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, types.Response{
		Status:  200,
		Message: "success",
		Data:    user,
	})
}

// @Tags Users
// @Summary Update user information
// @Description Update information about the currently logged in user
// @Accept  json
// @Produce  json
// @Param   user    body    types.RequestUpdateUser   true    "user info to update"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Security Bearer
// @Router /users/me [put]
func (controller *UserController) UpdateUserInfo(c *gin.Context) {
	var input types.RequestUpdateUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	user, err := controller.UserService.UpdateUserInfo(c, input)
	if err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, types.Response{
		Status:  200,
		Message: "success",
		Data:    user,
	})
}

// @Tags Users
// @Summary Delete user
// @Description Delete the currently logged in user
// @Accept  json
// @Produce  json
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Security Bearer
// @Router /users/me [delete]
func (controller *UserController) DeleteUser(c *gin.Context) {
	err := controller.UserService.DeleteUser(c)
	if err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: err.Error(),
			Data:    "fail",
		})
		return
	}

	c.JSON(200, types.Response{
		Status:  200,
		Message: "success",
		Data:    "OK",
	})
}
