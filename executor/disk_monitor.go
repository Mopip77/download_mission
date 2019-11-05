package executor

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type DiskMonitor struct {
	FullVolume  string  `json:"full_volume"`
	UsedVolume  string  `json:"used_volume"`
	UsedPercent float64 `json:"used_percent"`
}

var (
	G_DiskMonitor DiskMonitor
)

func (monitor *DiskMonitor) Run() {
	go func() {
		var (
			res []string
			f   float64
		)

		for {
			dir, _ := os.Getwd()
			scriptPath := path.Join(dir, "script", "for_this_proj", "get_df.sh")
			cmd := exec.Command(scriptPath)
			bytes, e := cmd.Output()
			if e != nil {
				log.Println(e)
				goto ERR
			}

			res = strings.Split(string(bytes), "\n")
			f, e = strconv.ParseFloat(strings.TrimSpace(res[0]), 64)
			if e != nil {
				log.Println(e)
				goto ERR
			}
			monitor.UsedPercent = f
			monitor.FullVolume = strings.TrimSpace(res[1])
			monitor.UsedVolume = strings.TrimSpace(res[2])

		ERR:
			time.Sleep(5 * time.Second)
		}
	}()
}
