package test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"onedrive/cache"
	"onedrive/executor"
	"os/exec"
	"testing"
	"time"
)

func Test_cmd(t *testing.T) {
	command := exec.Command("bash", "/Users/mopip77/project/go/onedrive/test/a.sh")
	go command.Run()

	time.Sleep(5 * time.Second)
	fmt.Println(command.ProcessState)
	stdoutPipe, _ := command.StdoutPipe()

	reader := bufio.NewReader(stdoutPipe)
	readString, _ := reader.ReadString('\n')
	fmt.Println(readString)
	//stdoutPipe, _ := command.StdoutPipe()
	//reader := bufio.NewReader(stdoutPipe)
	//s, _ := reader.ReadString('\n')
	//fmt.Println(s)
}

func Test_wo(t *testing.T) {
	cache.Redis()
	scan := cache.RedisClient.Scan(0, "/mission/*", 10)
	keys, cursor, _ := scan.Result()

	fmt.Println(keys)
	fmt.Println(cursor)
}

func Test_jsoon(t *testing.T) {
	str := "[{\"FolderPath\":\"/Users/mopip77/1572771567\",\"Urls\":[\"www.baidu.com\",\"www.google.com\"],\"DownloadFiles\":null,\"StartTime\":\"2019-11-03T16:59:27.41433+08:00\",\"DownloadIdx\":0,\"State\":\"dead\"},{\"FolderPath\":\"/Users/mopip77/1572771392\",\"Urls\":[\"www.baidu.com\",\"www.google.com\"],\"DownloadFiles\":null,\"StartTime\":\"2019-11-03T16:56:32.735667+08:00\",\"DownloadIdx\":0,\"State\":\"dead\"}]"
	var tt []executor.Mission
	json.Unmarshal([]byte(str), &tt)
	fmt.Println(tt)
}

func Test_json(t *testing.T) {
	var tt []executor.Mission
	tt = append(tt, *executor.NewMission([]string{"a", "b"}))
	tt = append(tt, *executor.NewMission([]string{"c", "d"}))
	bytes, _ := json.Marshal(tt)
	fmt.Println(string(bytes))
}

func Test_ls(t *testing.T) {
	cache.Redis()
	result := cache.RedisClient.Get("/mission/1572852826")
	if result.Val() != "" {
		var mission executor.Mission
		e := json.Unmarshal([]byte(result.Val()), &mission)
		fmt.Println(e)
		fmt.Println(mission)
	}
}
