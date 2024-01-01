package controllers

import (
	"gdsc/baro/types"

	"github.com/gin-gonic/gin"
)

type HealthCheckController struct{}

// @Tags HealthCheck
// @Summary HealthCheck
// @Description HealthCheck
// @Accept  json
// @Produce  json
// @Success 200 {object} types.Response
// @Router /health [get]
func (h HealthCheckController) HealthCheck(c *gin.Context) {
	c.JSON(200, types.Response{
		Status:  200,
		Message: "서버 정상 작동 중",
		Data:    "OK",
	})
}
