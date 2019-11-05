package executor

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CommandLs(args ...string) []string {
	cmd := exec.Command("ls", args...)
	bytes, e := cmd.Output()
	if e != nil {
		return []string{""}
	}
	return strings.Split(strings.TrimSpace(string(bytes)), "\n")
}

// 由于程序中有进度条显示，所以输出使用\r来重置进度条，但是如果不处理就会有很多进度条被读出来
// 该函数处理\r(backslash r)，通过字节流的遍历模拟\r的重置
func ReadFileHandleBackslashR(filePath string) (string, error) {
	bytes, e := readFileBytes(filePath)
	if e != nil {
		return "", e
	}
	var newBytes = make([]byte, len(bytes))
	var curLineIdx = 0
	var cursor = 0
	for _, byte := range bytes {
		if byte == '\r' {
			cursor = curLineIdx
		} else {
			newBytes[cursor] = byte
			cursor++
			if byte == '\n' {
				curLineIdx = cursor
			}
		}
	}
	return string(newBytes[:cursor]), nil
}

func ReadFile(filePath string) (string, error) {
	bytes, e := readFileBytes(filePath)
	return string(bytes), e
}

func readFileBytes(filePath string) ([]byte, error)  {
	file, e := os.Open(filePath)
	defer file.Close()
	if e != nil {
		log.Println(e)
		return nil, e
	}
	bytes, e := ioutil.ReadAll(file)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	return bytes, nil
}
