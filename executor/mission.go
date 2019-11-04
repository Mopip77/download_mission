package executor

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"io"
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
	// command
	cmd *exec.Cmd
	// 用于中断任务的函数
	cancelFunc func()
	// 下载url的列表文件
	urlFilePath string
	// 用于检测下载文件的chan
	ch chan bool
	// stdoutpipe 这个需要在进程执行前创建，所以必须先保存下来
	stdoutPipe io.ReadCloser
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
	folder := path.Join(os.Getenv("HOME"), nowUnix)
	urlFilePath := folder + ".txt"
	urlFile, err := os.Create(urlFilePath)
	if err != nil {
		return nil
	}
	defer urlFile.Close()
	_, err = urlFile.Write([]byte(strings.Join(urls, "\n")))
	if err != nil {
		return nil
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	command := exec.CommandContext(ctx, "bash", os.Getenv("DL_SCRIPT"), nowUnix, urlFilePath)
	mission := Mission{
		FolderPath:     folder,
		Urls:           urls,
		StartTime:      now,
		StartTimeStamp: now.Unix(),
		State:          INIT,
		cmd:            command,
		cancelFunc:     cancelFunc,
		urlFilePath:    urlFilePath,
		ch:             make(chan bool),
	}
	readCloser, err := command.StdoutPipe()
	if err == nil {
		mission.stdoutPipe = readCloser
	}
	return &mission
}

func (mission *Mission) Start() {
	log.Println("mission START:", mission.StartTimeStamp)
	go func() {
		defer close(mission.ch)
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
				if dlFiles := CommandLs(mission.FolderPath); len(dlFiles) >= len(mission.DownloadFiles) {
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
	if mission.stdoutPipe == nil {
		return "获取stdou出错..."
	}
	// ProcessState == nil 为程序正在运行，!= nil即程序运行完成，并且如果程序运行完成就不能读取stdout输出了
	if mission.cmd.ProcessState != nil {
		return "Process finished..."
	}

	reader := bufio.NewReader(mission.stdoutPipe)
	readString, e := reader.ReadString('\n')
	if e != nil {
		return ""
	}

	return readString
}

func (mission *Mission) MissionKeyOnRedis() string {
	return os.Getenv("MISSION_PREFIX") + strconv.Itoa(int(mission.StartTimeStamp))
}

func (mission *Mission) removeDownloadFile() {
	log.Println(mission.StartTimeStamp, " is ", mission.State, " remove files...")
	os.Remove(mission.urlFilePath)
	os.RemoveAll(mission.FolderPath)
}
