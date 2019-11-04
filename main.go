package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"onedrive/conf"
	"onedrive/server"
)

func main() {
	conf.Init()

	fmt.Println("dl api running...")
	gin.SetMode(gin.ReleaseMode)
	route:=server.InitRoute()
	route.Run(":3000")
}
