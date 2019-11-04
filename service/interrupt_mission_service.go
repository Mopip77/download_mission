package service

import (
	"onedrive/executor"
	"onedrive/serializer"
)

type InterruptMissionService struct {
	TimeStamp int64 `form:"time_stamp" json:"time_stamp" binding:"required"`
}

func (service *InterruptMissionService) Interrupt() serializer.Response {
	found := false
	for _, mission := range executor.G_Executor.Missions {
		if mission.StartTimeStamp == service.TimeStamp {
			mission.Interrupt()
			found = true
			break
		}
	}

	if found {
		return serializer.Response{
			Status: 200,
			Msg:    "任务终止成功",
		}
	} else {
		return serializer.Response{
			Status: 1001,
			Msg:    "未找到该时间对应任务",
		}
	}
}
