package executor

import (
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
