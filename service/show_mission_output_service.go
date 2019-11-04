package service

import (
	"fmt"
	"onedrive/executor"
	"onedrive/serializer"
)

type ShowMissionOutputService struct {
	TimeStamp int64 `form:"time_stamp" json:"time_stamp" binding:"required"`
}

func (service *ShowMissionOutputService) Output() serializer.Response {
	fmt.Println("ts:", service.TimeStamp)
	found := false
	output := ""
	// 先在正在执行的任务中查询
	for _, mission := range executor.G_Executor.Missions {
		if mission.StartTimeStamp == service.TimeStamp {
			output = mission.Output()
			found = true
			break
		}
	}

	if found {
		return serializer.Response{
			Status: 200,
			Data:   output,
		}
	} else {
		return serializer.Response{
			Status: 1001,
			Msg:    "未找到该时间对应任务",
		}
	}
}
