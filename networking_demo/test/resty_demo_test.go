package test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
)

func Test_restyDemo(t *testing.T) {
	checkUrl := "https://tls.peet.ws/api/all"
	//checkUrl := "https://www.linkedin.com/checkpoint/lg/login"
	//
	//networking_demo.RestyDemo(checkUrl)
	//networking_demo.RestyDemo(checkUrl)
	client := resty.New()

	resp, err := client.R().
		Get(checkUrl)
	if err != nil {
		fmt.Println("请求出错:", err)
		return
	}
	fmt.Println(resp.StatusCode())
	// 处理返回的响应
	fmt.Println(resp.String())
}
