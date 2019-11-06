package service

import (
	"onedrive/executor"
	"onedrive/serializer"
	"os"
	"path"
	"strconv"
)

type ShowMissionOutputService struct {
	TimeStamp int64 `form:"time_stamp" json:"time_stamp" binding:"required"`
}

func (service *ShowMissionOutputService) Output() serializer.Response {
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
			Data:   output,
		}
	} else {
		// 查文件
		missionLogPath := path.Join(os.Getenv("MISSION_LOG_PATH"), strconv.Itoa(int(service.TimeStamp))) + ".log"
		content, e := executor.ReadFileHandleBackslashR(missionLogPath)
		if e != nil {
			return serializer.Response{
				Status: 1002,
				Msg:    "读取任务输出出错",
				Error: e.Error(),
			}
		} else {
			return serializer.Response{
				Data:   content,
			}
		}
	}
}
