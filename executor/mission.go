package executor

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"onedrive/cache"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type MissionState string

const (
	INIT    MissionState = "init"
	RUNNING MissionState = "running"
	DEAD    MissionState = "dead"
	FINISH  MissionState = "finish"
)

// 每个you-get任务的结构
// 为了区分不同的下载任务，为每个任务分配一个独立的文件夹，并且使用时间戳命名
type Mission struct {
	// 下载文件夹
	FolderPath string `json:"folder_path"`
	// 需要下载的网址
	Urls []string `json:"urls"`
	// 下载的文件名
	DownloadFiles []string `json:"download_files"`
	// 开始时间
	StartTime time.Time `json:"start_time"`
	// 开始时间戳
	StartTimeStamp int64 `json:"start_time_stamp"`
	// 状态
	State MissionState `json:"state"`
	// 结束用时/秒, FINISH或DEAD状态的任务才有，其他默认为0
	MissionUsedTime int64 `json:"mission_used_time"`
	// command
	cmd *exec.Cmd
	// 用于中断任务的函数
	cancelFunc func()
	// 下载url的列表文件
	urlFilePath string
	// 日志文件
	logFile *os.File
	// 用于检测下载文件的chan
	ch chan bool
}

func CreateAndRun(urls []string) *Mission {
	mission := NewMission(urls)
	G_Executor.AddMission(mission)
	mission.RegisterKeyOnRedis()
	mission.UpdateOnRedis()
	mission.Start()
	return mission
}

func NewMission(urls []string) *Mission {
	now := time.Now()
	nowUnix := strconv.Itoa(int(now.Unix()))
	videoDLFolder := path.Join(os.Getenv("VIDEO_DL_PATH"), nowUnix)
	urlFilePath := videoDLFolder + ".txt"
	missionLogPath := path.Join(os.Getenv("MISSION_LOG_PATH"), nowUnix) + ".log"
	// 因为该程序的输出重定向到日志文件，所以该日志文件必须在任务执行结束后才能关闭，所以该文件流的关闭函数在start()后执行
	logFile, err := os.Create(missionLogPath)
	if err != nil {
		return nil
	}

	urlFile, err := os.Create(urlFilePath)
	defer urlFile.Close()
	if err != nil {
		return nil
	}

	_, err = urlFile.Write([]byte(strings.Join(urls, "\n")))
	if err != nil {
		return nil
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	command := exec.CommandContext(ctx, "bash", os.Getenv("DL_SCRIPT"), videoDLFolder, urlFilePath)
	command.Stdout = logFile
	mission := Mission{
		FolderPath:      videoDLFolder,
		Urls:            urls,
		StartTime:       now,
		StartTimeStamp:  now.Unix(),
		State:           INIT,
		MissionUsedTime: 0,
		cmd:             command,
		cancelFunc:      cancelFunc,
		urlFilePath:     urlFilePath,
		logFile:         logFile,
		ch:              make(chan bool),
	}
	return &mission
}

func (mission *Mission) Start() {
	log.Println("mission START:", mission.StartTimeStamp)
	go func() {
		defer close(mission.ch)
		defer mission.logFile.Close()
		defer mission.removeDownloadFile()
		mission.State = RUNNING
		mission.UpdateOnRedis()

		go func() {
			flag := true
			for flag {
				select {
				case <-mission.ch:
					flag = false
				case <-time.After(3 * time.Second):
					if stat, err := os.Stat(mission.FolderPath); err == nil && stat != nil {
						mission.DownloadFiles = CommandLs(mission.FolderPath)
						mission.UpdateOnRedis()
					}
				}
				//log.Println("[", mission.StartTime.Unix(), "]checking download files:", mission.DownloadFiles)
			}
			if stat, err := os.Stat(mission.FolderPath); err == nil && stat != nil {
				// 运行结束后还可以再查一次，可能在间隔中有漏掉的文件，也可能文件被删除了，所以需要先判断文件个数
				if dlFiles := CommandLs(mission.FolderPath); len(dlFiles) != 0 {
					mission.DownloadFiles = dlFiles
					mission.UpdateOnRedis()
				}
			}
			//log.Println("[", mission.StartTime.Unix(), "]checking goroutine finish...")
		}()
		e := mission.cmd.Run()
		if e != nil {
			mission.State = DEAD
			log.Println("mission DEAD:", mission.StartTimeStamp)
			log.Println("reason:", e)
		} else {
			mission.State = FINISH
			log.Println("mission FINISH:", mission.StartTimeStamp)
		}
		mission.MissionUsedTime = time.Now().Unix() - mission.StartTime.Unix()
		mission.ch <- true
		mission.UpdateOnRedis()
		G_Executor.DeleteMission(*mission)
	}()
}

func (mission *Mission) Interrupt() {
	log.Println("mission INTERRUPT:", mission.StartTimeStamp)
	mission.cancelFunc()
	mission.State = DEAD
	mission.UpdateOnRedis()
	G_Executor.DeleteMission(*mission)
}

func (mission *Mission) Retry() {
	log.Println("mission RETRY:", mission.StartTimeStamp)
	mission.Interrupt()
	// 清除原先redis中的mission
	cache.RedisClient.ZRem(os.Getenv("REDIS_ZSET_KEY"), mission.MissionKeyOnRedis())
	cache.RedisClient.Del(mission.MissionKeyOnRedis())

	// new Mission
	newMission := NewMission(mission.Urls)
	G_Executor.AddMission(newMission)
	newMission.RegisterKeyOnRedis()
	newMission.UpdateOnRedis()
	newMission.Start()
}

func (mission *Mission) UpdateOnRedis() {
	bytes, e := json.Marshal(mission)
	if e != nil {
		log.Println(e)
	}
	cache.RedisClient.GetSet(mission.MissionKeyOnRedis(), bytes)
}

// 由于需要从redis中拿任务，但是redis不按时间排序，所以维护一个zset，每次新建任务时插入key
func (mission *Mission) RegisterKeyOnRedis() {
	key := "mission_keys"
	cache.RedisClient.ZAdd(key, redis.Z{
		Score:  float64(mission.StartTimeStamp),
		Member: mission.MissionKeyOnRedis(),
	})
}

func (mission *Mission) Output() string {
	content, e := ReadFileHandleBackslashR(mission.logFile.Name())
	if e != nil {
		return ""
	} else {
		return content
	}
}

func (mission *Mission) MissionKeyOnRedis() string {
	return os.Getenv("MISSION_PREFIX") + strconv.Itoa(int(mission.StartTimeStamp))
}

// 任务终止或完成，log并删除下载文件
func (mission *Mission) removeDownloadFile() {
	log.Println(mission.StartTimeStamp, " is ", mission.State, " remove files...")
	os.Remove(mission.urlFilePath)
	os.RemoveAll(mission.FolderPath)
}
