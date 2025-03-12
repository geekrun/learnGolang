package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -llog

#include <stdlib.h>  // ✅ 解决 C.free() 问题
#include <android/log.h>

#define LOG_TAG "GoNDK"
#define LOGI(...) __android_log_print(ANDROID_LOG_INFO, LOG_TAG, __VA_ARGS__)

void logMessage(const char* message) {
    LOGI("%s", message);
}
*/
import "C"
import (
	"embed"
	"fmt"
	"os"
	"token_manager/cmd"
	"unsafe"
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

// Golang 调用 Android NDK 日志 API
func LogInfo(msg string) {
	msg = "GolangLog: " + msg
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg)) // ✅ 现在 C.free() 正确解析
	C.logMessage(cmsg)
}

// 保存脚本到 Android 设备本地
func extractScript(scriptName string) (string, error) {
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
	LogInfo("Hello from Golang using Android NDK!")
	// 提取并运行 `script1.sh`
	path1, err := extractScript("backup_gms_sms.sh")
	if err != nil {
		fmt.Println("Failed to extract script1.sh:", err)
		return
	}
	path2, err := extractScript("restore_gms_sms.sh")
	if err != nil {
		fmt.Println("Failed to extract script2.sh:", err)
		return
	}
	LogInfo("backup_gms_sms.sh " + path1)
	LogInfo("restore_gms_sms " + path2)
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
