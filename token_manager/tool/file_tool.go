package tool

import (
	"encoding/json"
	"os"
)

// ReadJSONFile 读取指定路径的 JSON 文件并返回 map[string]interface{}
func ReadJSONFile(filePath string) (map[string]interface{}, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 定义一个 map 来存储解析后的数据
	var result map[string]interface{}

	// 解析 JSON 数据
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
