package conf

import (
	"github.com/joho/godotenv"
	"log"
	"onedrive/cache"
	"onedrive/executor"
	"os"
	"path"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	godotenv.Load()

	// 指定下载脚本路径
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	os.Setenv("DL_SCRIPT", path.Join(curDir, "script", "deploy", "youget-dl.sh"))

	// 新建文件下载和日志文件夹
	os.Setenv("VIDEO_DL_PATH", path.Join(os.Getenv("HOME"), "video"))
	os.Setenv("MISSION_LOG_PATH", path.Join(os.Getenv("HOME"), "log"))
	os.MkdirAll(os.Getenv("VIDEO_DL_PATH"), os.ModePerm)
	os.MkdirAll(os.Getenv("MISSION_LOG_PATH"), os.ModePerm)

	// 连接数据库
	cache.Redis()

	// 启动disk usage检测
	executor.G_DiskMonitor.Run()
}
