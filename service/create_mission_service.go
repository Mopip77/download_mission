package service

import (
	"onedrive/executor"
	"onedrive/serializer"
)

type CreateMissionService struct {
	Urls []string `form:urls`
}

func (service *CreateMissionService) Create() serializer.Response {
	if len(service.Urls) == 0 {
		return serializer.Response{
			Status: 1000,
			Msg:    "下载视频列表为空",
		}
	}
	mission := executor.CreateAndRun(service.Urls)

	return serializer.Response{
		Status: 200,
		Data:   mission,
	}
}
