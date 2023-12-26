package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckController struct{}

func (h HealthCheckController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
