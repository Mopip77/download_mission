package api

import (
	"github.com/gin-gonic/gin"
	"onedrive/service"
)

// CreateMission
func CreateMission(c *gin.Context) {
	service := service.CreateMissionService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Create()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// ListMission
func ListMission(c *gin.Context) {
	service := service.ListMissionService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.List()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// InterruptMission
func InterruptMission(c *gin.Context) {
	service := service.InterruptMissionService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Interrupt()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// RerunMission
func RerunMission(c *gin.Context) {
	service := service.RerunMissionService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Rerun()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// ShowMissionOutput
func ShowMissionOutput(c *gin.Context) {
	service := service.ShowMissionOutputService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Output()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}