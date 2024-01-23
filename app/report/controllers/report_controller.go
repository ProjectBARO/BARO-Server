package controllers

import (
	"gdsc/baro/app/report/services"
	"gdsc/baro/app/report/types"
	"gdsc/baro/global"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	ReportService services.ReportServiceInterface
}

func NewReportController(reportService services.ReportServiceInterface) *ReportController {
	return &ReportController{
		ReportService: reportService,
	}
}

// @Tags Reports
// @Summary 자세 추정 요청
// @Description 자세를 추정합니다. (동영상 URL을 입력받아 자세를 추정합니다.)
// @Accept  json
// @Produce  json
// @Param   video_url    body    types.RequestAnalysis   true    "URL, 알림 횟수 등"
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /analysis [post]
func (controller *ReportController) Analysis(c *gin.Context) {
	var input types.RequestAnalysis
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

	response, err := controller.ReportService.Analysis(c, input)
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
// @Router /analysis [get]
func (controller *ReportController) GetAnalysis(c *gin.Context) {
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
// @Summary 자세 추정 결과 id로 조회
// @Description 보고서 id로 자세 추정 결과를 조회합니다. (요약으로 먼저 보고서 id 조회하고 사용자가 그걸 누르면 이걸 사용하기)
// @Accept  json
// @Produce  json
// @Param   id    path    int   true    "자세 추정 결과 id"
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /analysis/{id} [get]
func (controller *ReportController) GetAnalysisById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: "Invalid ID",
		})
		return
	}

	response, err := controller.ReportService.FindById(c, uint(id))
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
// @Summary 자세 추정 결과 월별 요약 조회
// @Description 로그인한 사용자의 자세 추정 결과를 월별로 요약하여 조회합니다. (캘린더 점 찍는 용도로 사용)
// @Accept  json
// @Produce  json
// @Param   ym  query    string   true    "조회할 년월 (YYYYMM) 예시: 202401 (2024년 1월)"
// @Success 200 {object} global.Response
// @Failure 400 {object} global.Response
// @Security Bearer
// @Router /analysis/summary [get]
func (controller *ReportController) GetAnalysisSummary(c *gin.Context) {
	yearAndMonth := c.Query("ym")

	response, err := controller.ReportService.FindReportSummaryByMonth(c, yearAndMonth)
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
// @Router /analysis/all [get]
func (controller *ReportController) GetAnalyzes(c *gin.Context) {
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
