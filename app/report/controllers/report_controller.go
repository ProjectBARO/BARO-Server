package controllers

import (
	"gdsc/baro/app/report/services"
	"gdsc/baro/app/report/types"
	"gdsc/baro/global"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	ReportService *services.ReportService
}

func NewReportController(reportService *services.ReportService) *ReportController {
	return &ReportController{
		ReportService: reportService,
	}
}

// @Tags Reports
// @Summary 자세 추정 요청
// @Description 자세를 추정합니다. (동영상 URL을 입력받아 자세를 추정합니다.)
// @Accept  json
// @Produce  json
// @Param   video_url    body    types.RequestPredict   true    "동영상 URL"
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /predict [post]
func (controller *ReportController) Predict(c *gin.Context) {
	var input types.RequestPredict
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

	response, err := controller.ReportService.Predict(c, input)
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
		Data:    response,
	})
}

// @Tags Reports
// @Summary 자세 추정 결과 조회
// @Description 로그인한 사용자의 자세 추정 결과를 조회합니다.
// @Accept  json
// @Produce  json
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /predict [get]
func (controller *ReportController) GetPredict(c *gin.Context) {
	response, err := controller.ReportService.FindReportByCurrentUser(c)
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
		Data:    response,
	})
}

// @Tags Reports
// @Summary 자세 추정 결과 전체 조회 (테스트용)
// @Description 자세 추정 결과를 조회합니다.
// @Accept  json
// @Produce  json
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /predict/all [get]
func (controller *ReportController) GetPredicts(c *gin.Context) {
	response, err := controller.ReportService.FindAll()
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
		Data:    response,
	})
}
