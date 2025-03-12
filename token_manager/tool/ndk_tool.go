package tool

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
	"os"
	"os/exec"
	"unsafe"
)

// RunShellScript 执行脚本
func RunShellScript(scriptName string, tokenId string) error {
	scriptPath := "/data/local/tmp/" + scriptName
	cmd := exec.Command("/system/bin/sh", scriptPath, tokenId)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// LogInfo Golang 调用 Android NDK 日志 API
func LogInfo(msg string) {
	msg = "GolangLog: " + msg
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg)) // ✅ 现在 C.free() 正确解析
	C.logMessage(cmsg)
}
