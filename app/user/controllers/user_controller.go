package controllers

import (
	"gdsc/baro/app/user/services"
	"gdsc/baro/app/user/types"
	"gdsc/baro/global"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserServiceInterface
}

func NewUserController(userService services.UserServiceInterface) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// @Tags Users
// @Summary 로그인 (첫 로그인 시 회원가입)
// @Description 토큰을 반환합니다. (첫 로그인 시 회원가입이 진행 후 토큰을 반환합니다.)
// @Accept  json
// @Produce  json
// @Param   user    body    types.RequestCreateUser   true    "사용자 정보"
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Router /login [post]
func (controller *UserController) LoginOrRegisterUser(c *gin.Context) {
	var input types.RequestCreateUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	token, err := controller.UserService.Login(input)
	if err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, global.Response{
		Status:  200,
		Message: "success",
		Data:    token,
	})
}

// @Tags Users
// @Summary FCM 토큰 업데이트
// @Description FCM 토큰을 업데이트합니다.
// @Accept  json
// @Produce  json
// @Param   fcm_token    body    types.RequestUpdateFcmToken   true    "FCM 토큰"
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /users/fcm-token [put]
func (controller *UserController) UpdateFcmToken(c *gin.Context) {
	var input types.RequestUpdateFcmToken
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	err := controller.UserService.UpdateFcmToken(c, input)
	if err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
			Data:    "fail",
		})
		return
	}

	c.JSON(200, global.Response{
		Status:  200,
		Message: "success",
		Data:    "OK",
	})
}

// @Tags Users
// @Summary 내 정보 조회
// @Description 현재 로그인한 사용자의 정보를 조회합니다.
// @Accept  json
// @Produce  json
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /users/me [get]
func (controller *UserController) GetUserInfo(c *gin.Context) {
	user, err := controller.UserService.GetUserInfo(c)
	if err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, global.Response{
		Status:  200,
		Message: "success",
		Data:    user,
	})
}

// @Tags Users
// @Summary 내 정보 수정
// @Description 현재 로그인한 사용자의 정보를 수정합니다.
// @Accept  json
// @Produce  json
// @Param   user    body    types.RequestUpdateUser   true    "수정할 사용자 정보"
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /users/me [put]
func (controller *UserController) UpdateUserInfo(c *gin.Context) {
	var input types.RequestUpdateUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	user, err := controller.UserService.UpdateUserInfo(c, input)
	if err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, global.Response{
		Status:  200,
		Message: "success",
		Data:    user,
	})
}

// @Tags Users
// @Summary 내 정보 삭제 (회원 탈퇴)
// @Description 현재 로그인한 사용자의 정보를 삭제합니다.
// @Accept  json
// @Produce  json
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /users/me [delete]
func (controller *UserController) DeleteUser(c *gin.Context) {
	err := controller.UserService.DeleteUser(c)
	if err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
			Data:    "fail",
		})
		return
	}

	c.JSON(200, global.Response{
		Status:  200,
		Message: "success",
		Data:    "OK",
	})
}
