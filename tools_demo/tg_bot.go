package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

var rdb *redis.Client

func redisInit() {
	//初始化redis，连接地址和端口，密码，数据库名称
	rdb = redis.NewClient(&redis.Options{
		Addr:     "35.187.225.32:6379",
		Password: "aje25!*djDK2S",
		DB:       2,
	})
}

type TimeoutConfig struct {
	MaxIntervalUpdatePeriodTime  int64 `json:"max_interval_update_period_time"`
	MaxIntervalSuccessPeriodTime int64 `json:"max_interval_success_period_time"`
}

type MonitorKeyLabel struct {
	LastUpdateTime  int64 `json:"last_update_time"`
	LastSuccessTime int64 `json:"last_success_time"`
	LastInstallTime int64 `json:"last_install_time"`
}

type AppData struct {
	AppName         string          `json:"app_name"`
	MachineLabel    string          `json:"machine_label"`
	DeviceLabel     string          `json:"device_label"`
	UploadTime      int64           `json:"upload_time"`
	TimeoutConfig   TimeoutConfig   `json:"timeout_config"`
	MonitorKeyLabel MonitorKeyLabel `json:"monitor_key_label"`
}

func checkMessageFromRedis(redisQueueName string) string {
	resultStr := ""
	ctx := context.Background()
	// 从 Redis 哈希队列中获取所有键值信息
	result, err := rdb.HGetAll(ctx, redisQueueName).Result()
	if err != nil {
		fmt.Println("Failed to get values from Redis:", err)
		return resultStr
	}

	// 遍历键值信息并反序列化值信息
	for key, value := range result {
		// 反序列化值信息到结构体
		var appdata AppData
		err := json.Unmarshal([]byte(value), &appdata)
		if err != nil {
			fmt.Println("Failed to deserialize value:", err)
			continue
		}
		curTimestamp := time.Now().Unix()
		bytes, _ := json.Marshal(appdata)

		errorMsg := ""
		if curTimestamp-appdata.UploadTime > appdata.TimeoutConfig.MaxIntervalUpdatePeriodTime {
			uploadTimeString := time.Unix(appdata.UploadTime, 0).Format("2006-01-02 15:04:05")
			errorMsg = fmt.Sprintf("服务上报时间超时,UploadTime %s ，超出最大间隔 %d秒", uploadTimeString, appdata.TimeoutConfig.MaxIntervalUpdatePeriodTime)
		} else if curTimestamp-appdata.MonitorKeyLabel.LastUpdateTime > appdata.TimeoutConfig.MaxIntervalUpdatePeriodTime {
			updateTimeString := time.Unix(appdata.MonitorKeyLabel.LastUpdateTime, 0).Format("2006-01-02 15:04:05")
			errorMsg = fmt.Sprintf("apk上报时间超时,LastUpdateTime %s 超出最大间隔 %d秒", updateTimeString, appdata.TimeoutConfig.MaxIntervalUpdatePeriodTime)
		} else if curTimestamp-appdata.MonitorKeyLabel.LastSuccessTime > appdata.TimeoutConfig.MaxIntervalSuccessPeriodTime {
			uploadTimeString := time.Unix(appdata.MonitorKeyLabel.LastSuccessTime, 0).Format("2006-01-02 15:04:05")
			errorMsg = fmt.Sprintf("上报服务出现问题，LastSuccessTime %s , 超出最大间隔 %d秒", uploadTimeString, appdata.TimeoutConfig.MaxIntervalSuccessPeriodTime)
		}

		if errorMsg != "" {
			//    machine_label app_name  device_label
			//    应用 {} 机器 {} 设备 {} 出现问题，    上报更新时间最长 1小时， 成功最长间隔 2小时， 请即使关注
			curTimeString := time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05")
			resultStr = fmt.Sprintf("当前时间 %s \n 设备 %s,出现问题 \n 怀疑是,%s  \n 详细信息如下: %s", curTimeString, key, errorMsg, bytes)
			fmt.Println(key, "结果异常", resultStr)
			sendMessage(resultStr)

		} else {
			fmt.Println(key, "结果正常", string(bytes))
		}

	}

	return "ok"

}

func sendMessage(text string) {
	// 替换为您的 Bot API 密钥
	botAPIKey := "6213700454:AAEwTZcM9prdgtq3NkxwnoYUY7Vxar7HKSw"
	YourChatId := -1001576693542

	// 创建 Bot API 客户端
	bot, err := tgbotapi.NewBotAPI(botAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	// 设置 Bot 的 Debug 模式（可选）
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// 设置要发送消息的 chat ID
	chatID := int64(YourChatId)

	// 创建要发送的消息
	message := tgbotapi.NewMessage(chatID, text)

	// 发送消息
	_, err = bot.Send(message)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	redisInit()
	result := checkMessageFromRedis("monitor_map_queue")
	log.Println("checkMessageFromRedis result", result)
	//if result != "ok" {
	//	sendMessage("Golang Telegram Bot 出现意外")
	//}

}
