package controllers

import (
	"gdsc/baro/app/video/services"
	"gdsc/baro/global"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	VideoService *services.VideoService
}

func NewVideoController(videoService *services.VideoService) *VideoController {
	return &VideoController{
		VideoService: videoService,
	}
}

// @Tags Videos
// @Summary 유튜브 영상 목록 조회 (현재 키워드: 거북이, 스트레칭)
// @Description 설정한 키워드에 맞는 유튜브 영상 목록 50개를 조회합니다.
// @Accept  json
// @Produce  json
// @Security Bearer
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
