package server

import (
	"github.com/gin-gonic/gin"
	"onedrive/api"
)

func InitRoute() *gin.Engine {
	r := gin.Default()

	g := r.Group("/api")
	{
		g.POST("/mission", api.CreateMission)
		g.GET("/missions", api.ListMission)
		g.POST("/mission_cancel", api.InterruptMission)
		g.POST("/mission_rerun", api.RerunMission)
		// 由于程序原本使用了进度条覆盖的方法，所以stdout会有很多无用的进度条内容，所以先不使用该接口
		g.GET("/mission_output", api.ShowMissionOutput)
		g.GET("/disk_usage", api.VpsDiskUsage)
	}
	return r
}
