package test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"onedrive/cache"
	"onedrive/executor"
	"os/exec"
	"strconv"
	"strings"
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

func Test_moni(t *testing.T) {
	cmd := exec.Command("/Users/mopip77/project/go/onedrive/script/for_this_proj/get_df.sh")
	bytes, e := cmd.Output()
	if e != nil {
		log.Fatal(e)
	}

	res := strings.Split(string(bytes), "\n")
	f, e := strconv.ParseFloat(strings.TrimSpace(res[0]), 64)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(f)
	fmt.Println(strings.TrimSpace(res[1]))
	fmt.Println(strings.TrimSpace(res[2]))
}

func Test_time(t *testing.T) {
	before := time.Now()
	time.Sleep(5 * time.Second)
	now := time.Now()
	fmt.Println("before,", before.Unix())
	fmt.Println("now,", now.Unix())
	fmt.Println(now.Unix() - before.Unix())
}

func Test_output(t *testing.T) {
	var a []byte
	fmt.Println(string(a))
}

func Test_dl(t *testing.T) {
	scriptPath := "/Users/mopip77/project/go/onedrive/script/deploy/youget-dl.sh"
	cmd := exec.Command(scriptPath, "/Users/mopip77/Downloads/tt/cc", "/Users/mopip77/Downloads/tt/b")
	bytes, e := cmd.Output()
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Println(string(bytes))
}












