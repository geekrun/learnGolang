package main

import (
	"embed"
	"fmt"
	"os"
	"token_manager/cmd"
	"token_manager/tool"
)

/*

global_env()
token_info()
auto_clear_token()
backup_token()
restore_token(string number)

*/

//go:embed scripts/*
var scripts embed.FS

// ExtractScript 保存脚本到 Android 设备本地
func ExtractScript(scriptName string) (string, error) {
	data, err := scripts.ReadFile("scripts/" + scriptName)
	if err != nil {
		return "", err
	}

	// Android 设备存放脚本的路径
	scriptPath := "/data/local/tmp/" + scriptName
	err = os.WriteFile(scriptPath, data, 0755)
	if err != nil {
		return "", err
	}

	return scriptPath, nil
}

func main() {
	tool.LogInfo("Hello from Golang using Android NDK!")
	// 提取并运行 `script1.sh`
	path1, err := ExtractScript("backup_gms_sms_v2.sh")
	if err != nil {
		fmt.Println("Failed to extract script1.sh:", err)
		return
	}
	path2, err := ExtractScript("restore_gms_sms_v2.sh")
	if err != nil {
		fmt.Println("Failed to extract script2.sh:", err)
		return
	}
	tool.LogInfo("backup_gms_sms_v2.sh " + path1)
	tool.LogInfo("restore_gms_sms " + path2)
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
