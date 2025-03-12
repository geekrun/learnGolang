package tool

import "C"
import (
	"os"
	"os/exec"
)

// 执行脚本
func RunShellScript(scriptName string, tokenId string) error {
	scriptPath := "/data/local/tmp/" + scriptName
	cmd := exec.Command("/system/bin/sh", scriptPath, tokenId)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
