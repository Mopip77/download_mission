package service

import (
	"onedrive/executor"
	"onedrive/serializer"
)

type VpsDiskUsageService struct {

}

func (service *VpsDiskUsageService) Usage() serializer.Response {
	return serializer.Response{
		Data:   executor.G_DiskMonitor,
	}
}