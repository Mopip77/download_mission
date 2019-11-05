package api

import (
	"github.com/gin-gonic/gin"
	"onedrive/service"
)

// VpsDiskUsage
func VpsDiskUsage(c *gin.Context) {
	service := service.VpsDiskUsageService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Usage()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}