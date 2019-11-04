package service

import (
	"onedrive/cache"
	"onedrive/executor"
	"onedrive/serializer"
	"os"
	"time"
)

type ListMissionService struct {
	Showall bool `form:"showall" json:"showall"`
	Page    int  `form:"page" json:"page"`
	Size    int  `form:"size" json:"size"`
}

func (service *ListMissionService) correctEdge() {
	if service.Page < 0 {
		service.Page = 0
	}

	if service.Size <= 0 {
		service.Size = 10
	}
}

func (service *ListMissionService) List() serializer.Response {
	hasResult := false
	var res interface{}
	service.correctEdge()
	if service.Showall == true {
		// 展示所有mission，包括执行完成或执行失败的，按照时间顺序
		zRange := cache.RedisClient.ZRevRange(os.Getenv("REDIS_ZSET_KEY"), 0, time.Now().Unix())
		keys := zRange.Val()
		var offset int
		var rightBound int
		resultStr := "["
		// 如果越界就直接返回空
		if offset = service.Page * service.Size; offset < len(keys) {
			if rightBound = offset + service.Size; rightBound > len(keys) {
				rightBound = len(keys)
			}
			for idx, key := range keys[offset:rightBound] {
				mission := cache.RedisClient.Get(key).Val()
				resultStr += mission
				//fmt.Println(mission)
				if idx != rightBound - 1 {
					resultStr += ","
				}
			}
			resultStr += "]"
			hasResult = true
			res = resultStr
		}
	} else {
		// 只展示正在执行的mission
		missions := executor.G_Executor.Missions
		if missions != nil {
			hasResult = true
			res = missions
		}
	}
	if !hasResult {
		res = "[]"
	}
	return serializer.Response{
		Status: 200,
		Data:   res,
	}
}
