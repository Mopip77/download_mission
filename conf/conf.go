package conf

import (
	"github.com/joho/godotenv"
	"onedrive/cache"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	godotenv.Load()

	// 连接数据库
	cache.Redis()
}
