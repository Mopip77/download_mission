package executor

import (
	"sync"
)

var (
	G_Executor Executor
)

type Executor struct {
	sync.RWMutex
	Missions []*Mission
}

func (executor *Executor) GetMission(startTime int64) int {
	executor.RLock()
	defer executor.RUnlock()

	return executor.findWithoutLock(startTime)
}

func (executor *Executor) AddMission(mission *Mission) {
	executor.Lock()
	executor.Missions = append(executor.Missions, mission)
	executor.Unlock()
}

func (executor *Executor) findWithoutLock(startTime int64) int {
	for idx, m := range executor.Missions {
		if startTime == m.StartTimeStamp {
			return idx
		}
	}
	return -1
}

func (executor *Executor) DeleteMission(mission Mission) bool {
	executor.Lock()
	defer executor.Unlock()
	idx := executor.findWithoutLock(mission.StartTimeStamp)
	if idx >= 0 {
		executor.Missions = append(executor.Missions[:idx], executor.Missions[idx+1:]...)
		return true
	}
	return false
}

func (executor *Executor) DeleteMissionByStartTime(startTime int64) bool {
	executor.Lock()
	defer executor.Unlock()
	idx := executor.findWithoutLock(startTime)
	if idx >= 0 {
		executor.Missions = append(executor.Missions[:idx], executor.Missions[idx+1:]...)
		return true
	}
	return false
}
