package controllers

import (
	"gdsc/baro/app/video/services"
	"gdsc/baro/global"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	VideoService services.VideoServiceInterface
}

func NewVideoController(videoService services.VideoServiceInterface) *VideoController {
	return &VideoController{
		VideoService: videoService,
	}
}

// @Tags Videos
// @Summary 유튜브 영상 전체 목록 조회
// @Description 전체 유튜브 영상 목록을 조회합니다.
// @Accept  json
// @Produce  json
// @Success 200 {object} global.Response
// @Router /videos [get]
func (controller *VideoController) GetVideos(c *gin.Context) {
	videos, err := controller.VideoService.GetVideos(c)
	if err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
			Data:    "failed",
		})
		return
	}

	c.JSON(200, global.Response{
		Status:  200,
		Message: "success",
		Data:    videos,
	})
}

// @Tags Videos
// @Summary 키워드 검색 유튜브 영상 목록 조회
// @Description 설정한 키워드에 맞는 유튜브 영상 목록을 조회합니다.
// @Accept  json
// @Produce  json
// @Param   keyword     query    string     true        "검색할 키워드"
// @Success 200 {object} global.Response
// @Router /videos/category [get]
func (controller *VideoController) GetVideosByCategory(c *gin.Context) {
	videos, err := controller.VideoService.GetVideosByCategory(c)
	if err != nil {
		c.JSON(400, global.Response{
			Status:  400,
			Message: err.Error(),
			Data:    "failed",
		})
		return
	}

	c.JSON(200, global.Response{
		Status:  200,
		Message: "success",
		Data:    videos,
	})
}
