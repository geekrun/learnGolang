package tokenmanager

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"token_manager/tool"
)

// TokenManager 负责 Token 备份 & 恢复
type TokenManager struct {
	mu          sync.Mutex
	phoneNumber string
	versionCode string
}

// NewTokenManager 创建一个新的 TokenManager
func NewTokenManager(versionCode string) *TokenManager {

	if versionCode == "" {
		versionCode = "v2"
	}

	var phoneNumber = ""
	result, err := tool.ReadJSONFile("/data/local/tmp/fake_fingerprint.json")
	if err != nil {
		log.Fatal(err)
	}
	// 获取 sim_info 中的 phoneNumber
	if simInfo, ok := result["sim_info"].(map[string]interface{}); ok {
		if phoneNumber, ok := simInfo["phoneNumber"].(string); ok {
			fmt.Println("Phone Number:", phoneNumber)
		}
	}

	return &TokenManager{
		phoneNumber: phoneNumber,
		versionCode: versionCode,
	}
}

// BackupToken 备份 Token（压缩文件）
func (tm *TokenManager) BackupToken(phoneNumber string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if phoneNumber == "" {
		phoneNumber = tm.phoneNumber
	}
	err := tool.RunShellScript("backup_gms_sms.sh", phoneNumber)
	if err != nil {
		return err
	}
	return nil

}

// RestoreToken 通过 ID 恢复 Token（解压）
func (tm *TokenManager) RestoreToken(phoneNumber string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if phoneNumber == "" {
		return errors.New("备份号码必须存在")
	}
	err := tool.RunShellScript("restore_gms_sms.sh", phoneNumber)
	if err != nil {
		return err
	}
	fmt.Println("Token restored to:")
	return nil
}
