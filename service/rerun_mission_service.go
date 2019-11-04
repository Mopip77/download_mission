package service

import (
	"encoding/json"
	"onedrive/cache"
	"onedrive/executor"
	"onedrive/serializer"
	"os"
	"strconv"
)

type RerunMissionService struct {
	TimeStamp int64 `form:"time_stamp" json:"time_stamp" binding:"required"`
}

func (service *RerunMissionService) Rerun() serializer.Response {
	found := false
	// 先在正在执行的任务中查询
	for _, mission := range executor.G_Executor.Missions {
		if mission.StartTimeStamp == service.TimeStamp {
			mission.Retry()
			found = true
			break
		}
	}

	if found {
		return serializer.Response{
			Status: 200,
			Msg:    "任务重启成功",
		}
	} else {
		// 在redis中查询历史记录
		result := cache.RedisClient.Get(os.Getenv("MISSION_PREFIX") + strconv.Itoa(int(service.TimeStamp)))
		if result.Val() != "" {
			var mission executor.Mission
			e := json.Unmarshal([]byte(result.Val()), &mission)
			if e == nil {
				newMission := executor.CreateAndRun(mission.Urls)
				return serializer.Response{
					Status: 0,
					Data:   newMission,
				}
			}
		}
		return serializer.Response{
			Status: 1001,
			Msg:    "未找到该时间对应任务",
		}
	}
}


