package controllers

import (
	"gdsc/baro/services"
	"gdsc/baro/types"

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
// @Summary Get videos
// @Description Get all videos
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} types.Response
// @Router /videos [get]
func (controller *VideoController) GetVideos(c *gin.Context) {
	videos, err := controller.VideoService.GetVideos(c)
	if err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: err.Error(),
			Data:    "failed",
		})
		return
	}

	c.JSON(200, types.Response{
		Status:  200,
		Message: "success",
		Data:    videos,
	})
}
