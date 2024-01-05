package global

import (
	"github.com/gin-gonic/gin"
)

type HealthCheckController struct{}

// @Tags HealthCheck
// @Summary 서버 상태 확인
// @Description 서버가 정상 작동 중인지 확인합니다.
// @Accept  json
// @Produce  json
// @Success 200 {object} Response
// @Router /health [get]
func (h HealthCheckController) HealthCheck(c *gin.Context) {
	c.JSON(200, Response{
		Status:  200,
		Message: "서버 정상 작동 중",
		Data:    "OK",
	})
}
